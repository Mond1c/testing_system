package internal

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
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

func (ctx *CodeRunnerContext) makeExecutable() error {
	_, err := exec.Command("chmod", "+x", ctx.executablePath).Output()
	return err
}

func (ctx *CodeRunnerContext) compileCpp() error {
	_, err := exec.Command("g++", "-std=c++20", "-o", ctx.executablePath, ctx.filePath).Output()
	if err != nil {
		return err
	}
	return ctx.makeExecutable()
}

func (ctx *CodeRunnerContext) compileJava() error {
	panic("implement me")
}

func (ctx *CodeRunnerContext) compileGo() error {
	_, err := exec.Command("go", "build", "-o", ctx.executablePath, ctx.filePath).Output()
	if err != nil {
		return err
	}
	return ctx.makeExecutable()
}

// compileProgram compiles source code to executable file using giving CodeRunnerContext.
// Using specific compiler based on given language.
func (ctx *CodeRunnerContext) compileProgram() error {
	if _, err := os.Stat(ctx.filePath); os.IsNotExist(err) {
		return err
	}
	switch ctx.language {
	case "cpp":
		return ctx.compileCpp()
	case "go":
		return ctx.compileGo()
	default:
		return fmt.Errorf("unsupported language: %s", ctx.language)
	}
}

// compareOutput compares output with test case output.
func compareOutput(original, output string) TestResult {
	if original == output {
		return OK
	}
	return WA
}

// runTest runs test and return test result with giving CodeRunnerContext.
// test runs using executable with executablePath file.
func (ctx *CodeRunnerContext) runTest(directoryWithTests string, number int) (TestResult, error) {
	file, err := os.Open(directoryWithTests + "/" + strconv.FormatInt(int64(number), 10) + ".out")
	if err != nil {
		log.Fatal(err)
		return RE, err
	}
	defer file.Close()
	original := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		original += scanner.Text() + "\n"
	}
	inputFile, err := os.Open(directoryWithTests + "/" + strconv.FormatInt(int64(number), 10) + ".in")
	if err != nil {
		log.Fatal(err)
		return RE, err
	}
	defer inputFile.Close()
	cmd := exec.Command("./" + ctx.executablePath)
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
	err := os.Remove(ctx.executablePath)
	if err != nil {
		log.Fatal(err)
	}
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
