// Package pkg contains internal logic of the application.
package pkg

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// CodeRunnerContext is a struct that contains all necessary information for running tests.
type CodeRunnerContext struct {
	filePath       string
	language       string
	executablePath string
	results        []chan TestingResult
	threads        int
	failed         bool
	timeLimit      time.Duration
	compileFile    string
}

// NewCodeRunnerContext creates new COdeRunnerContext with giving parameters.
func NewCodeRunnerContext(filePath, executablePath, language string, compileFile string) *CodeRunnerContext {
	threads := 4
	results := make([]chan TestingResult, threads)
	for i := 0; i < threads; i++ {
		results[i] = make(chan TestingResult, 1)
	}
	return &CodeRunnerContext{
		filePath:       filePath,
		language:       language,
		executablePath: executablePath,
		results:        results,
		threads:        threads,
		timeLimit:      2.0,
		compileFile:    compileFile,
	}
}

// TestingResult is a struct that contains information about test result.
type TestingResult struct {
	Number int        `json:"number"`
	Result TestResult `json:"result"`
}

// String return string representation of TestingResult.
func (t *TestingResult) String() string {
	return fmt.Sprintf("Test with number %d: %s", t.Number, t.Result.GetString())
}

// compileProgram compiles source code to executable file using giving CodeRunnerContext.
// Using specific compiler based on given language.
func (ctx *CodeRunnerContext) compileProgram() error {
	if _, err := os.Stat(ctx.filePath); os.IsNotExist(err) {
		return err
	}
	_, err := exec.Command("sh", ctx.compileFile, ctx.filePath, ctx.executablePath).Output()
	log.Printf(
		"FileName: %v, path: %v, execPath: %v, err: %v",
		ctx.compileFile,
		ctx.filePath,
		ctx.executablePath,
		err,
	)
	return err
}

// compareOutput compares output with test case output.
func compareOutput(original, output string) TestResult {
	original = strings.Trim(original, " \t\n\r")
	output = strings.Trim(output, " \t\n\r")
	if original == output {
		return OK
	}
	return WA
}

// getExpectedOutput gets expected output from the file with the specified path
func (ctx *CodeRunnerContext) getExpectedOutput(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	original := ""
	for scanner.Scan() {
		original += scanner.Text() + "\n"
	}
	return original, nil
}

// getExecutableDirectoryAndFile gets directory and file name from executablePath
func (ctx *CodeRunnerContext) getExecutableDirectoryAndFile() (string, string) {
	path := strings.Split(ctx.executablePath, "/")
	directory := strings.Join(path[:len(path)-1], "/")
	file := path[len(path)-1]
	return directory, file
}

// runTest runs test and return test result with giving CodeRunnerContext.
// test runs using executable with executablePath file.
func (ctx *CodeRunnerContext) runTest(directoryWithTests string, number int) (TestResult, error) {
	path := fmt.Sprintf("%s/%d", directoryWithTests, number) // TODO: Think about windows separator
	original, err := ctx.getExpectedOutput(path + ".out")    // TODO: Think about performance with many users
	if err != nil {
		log.Fatal(err)
		return RE, err
	}
	inputFile, err := os.Open(path + ".in")
	if err != nil {
		log.Println(err)
		return RE, err
	}
	defer inputFile.Close()
	var cmd *exec.Cmd
	if ctx.language == "java" {
		dir, file := ctx.getExecutableDirectoryAndFile()
		cmd = exec.Command("java", file)
		cmd.Dir = dir
	} else {
		cmd = exec.Command("./" + ctx.executablePath)
	}
	type output struct {
		out []byte
		err error
	}
	ch := make(chan output)
	go func() {
		cmd.Stdin = inputFile
		out, err := cmd.CombinedOutput()
		ch <- output{out, err}
	}()
	select {
	case <-time.After(ctx.timeLimit * time.Second):
		err = cmd.Process.Kill()
		if err != nil {
			log.Print(err)
		}
		return TL, nil
	case x := <-ch:
		return compareOutput(original, string(x.out)), nil
	}
}

// removeExecutable removes file that creates after compilation
func (ctx *CodeRunnerContext) removeExecutable() {
	path := ctx.executablePath
	if ctx.language == "java" {
		path += ".class"
	}
	removeFile(path)
}

// runPartTests runs part of the tests with the specified start and end indexes.
func (ctx *CodeRunnerContext) runPartTests(directoryWithTests string, start, end, number int) {
	log.Println(start, end)
	for i := start; i < end; i++ {
		if ctx.failed {
			break
		}
		testResult, err := ctx.runTest(directoryWithTests, i)
		if err != nil || testResult != OK {
			ctx.failed = true
			log.Print(err)
			ctx.results[number] <- TestingResult{Number: i, Result: testResult}
			return
		}
	}
	ctx.results[number] <- TestingResult{Number: -1, Result: OK}
}

// Test tests program on given tests and returns result of testing
func (ctx *CodeRunnerContext) Test(
	directoryWithTests string,
	testsCount int,
) (TestingResult, error) {
	start := time.Now()
	err := ctx.compileProgram()
	defer ctx.removeExecutable()
	if err != nil {
		log.Println(err)
		return TestingResult{Number: -1, Result: CE}, err
	}
	step := testsCount / ctx.threads // maybe make constant
	if step < 1 {
		step = 1
	}
	var wg sync.WaitGroup
	for i := 0; i < testsCount; i += step {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			ctx.runPartTests(directoryWithTests, start, start+step, start/step)
		}(i)
	}
	wg.Wait()
	resultsCount := ctx.threads
	if testsCount < ctx.threads {
		resultsCount = 1
	}
	for i := 0; i < resultsCount; i++ {
		result := <-ctx.results[i]
		if result.Result != OK {
			return result, nil
		}
	}
	log.Printf("Time elapsed: %v", time.Since(start))
	return TestingResult{Number: -1, Result: OK}, nil
}
