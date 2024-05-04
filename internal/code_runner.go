package internal

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

type CodeRunnerContext struct {
	filePath       string
	language       string
	executablePath string
	results        []chan TestingResult
	threads        int
	failed         bool
}

func NewCodeRunnerContext(filePath, executablePath, language string) *CodeRunnerContext {
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
	}
}

type TestingResult struct {
	Number int        `json:"number"`
	Result TestResult `json:"result"`
	Err    error
}

func (t *TestingResult) GetString() string {
	return fmt.Sprintf("Test with number %d: %s", t.Number, t.Result.GetString())
}

// compileProgram compiles source code to executable file using giving CodeRunnerContext.
// Using specific compiler based on given language.
func (ctx *CodeRunnerContext) compileProgram() error {
	if _, err := os.Stat(ctx.filePath); os.IsNotExist(err) {
		return err
	}
	if fileName, ok := LangaugesConfig.CompileFiles[ctx.language]; ok {
		_, err := exec.Command("sh", fileName, ctx.filePath, ctx.executablePath).Output()
		log.Printf("FileName: %v, path: %v, execPath: %v, err: %v", fileName, ctx.filePath, ctx.executablePath, err)
		return err
	}
	return fmt.Errorf("unsupported language: %s", ctx.language)
}

// compareOutput compares output with test case output.
func compareOutput(original, output string) TestResult {
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

// runTest runs test and return test result with giving CodeRunnerContext.
// test runs using executable with executablePath file.
func (ctx *CodeRunnerContext) runTest(directoryWithTests string, number int) (TestResult, error) {
	path := fmt.Sprintf("%s/%d", directoryWithTests, number)
	original, err := ctx.getExpectedOutput(path + ".out")
	if err != nil {
		log.Fatal(err)
		return RE, err
	}
	inputFile, err := os.Open(path + ".in")
	if err != nil {
		log.Fatal(err)
		return RE, err
	}
	defer inputFile.Close()
	var cmd *exec.Cmd
	if ctx.language == "java" {
		cmd = exec.Command("java", ctx.executablePath)
	} else {
		cmd = exec.Command("./" + ctx.executablePath)
	}
	cmd.Stdin = inputFile
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("%s\n", output)
		return RE, err
	}

	return compareOutput(original, string(output)), nil
}

// removeExecutable removes file that creates after compilation
func (ctx *CodeRunnerContext) removeExecutable() {
	path := ctx.executablePath
	if ctx.language == "java" {
		path += ".class"
	}
	RemoveFile(path)
}

// runPartTests runs part of the tests with the specified start and end indexes.
func (ctx *CodeRunnerContext) runPartTests(directoryWithTests string, start, end, number int) {
	for i := start; i < end; i++ {
		if ctx.failed {
			break
		}
		testResult, err := ctx.runTest(directoryWithTests, i)
		if err != nil || testResult != OK {
			ctx.failed = true
			ctx.results[number] <- TestingResult{Number: i, Result: testResult, Err: err}
			return
		}
	}
	ctx.results[number] <- TestingResult{Number: -1, Result: OK, Err: nil}
}

// Test tests program on given tests and returns result of testing
func (ctx *CodeRunnerContext) Test(directoryWithTests string, testsCount int) (TestingResult, error) {
	start := time.Now()
	err := ctx.compileProgram()
	defer ctx.removeExecutable()
	if err != nil {
		log.Println(err)
		return TestingResult{Number: -1, Result: CE}, err
	}
	step := testsCount / ctx.threads // maybe make constant
	var wg sync.WaitGroup
	for i := 0; i < testsCount; i += step {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx.runPartTests(directoryWithTests, i, i+step, i/step)
		}()
	}
	wg.Wait()
	for i := 0; i < ctx.threads; i++ {
		result := <-ctx.results[i]
		if result.Result != OK || result.Err != nil {
			return result, result.Err
		}
	}
	log.Printf("Time elapsed: %v", time.Since(start))
	return TestingResult{Number: -1, Result: OK}, nil
}
