package internal

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
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
	output = strings.Trim(output, "\n")
	log.Println(original, output)
	if original == output {
		return OK
	}
	return WA
}

// runTest runs test and return test result with giving CodeRunnerContext.
// test runs using executable with executablePath file.
func (ctx *CodeRunnerContext) runTest(test *Test) (TestResult, error) {
	cmd := exec.Command("./" + ctx.executablePath)
	input, err := cmd.StdinPipe()
	if err != nil {
		return RE, err
	}

	go func() {
		defer func(input io.WriteCloser) {
			err := input.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(input)
		_, err := io.WriteString(input, test.input)
		if err != nil {
			log.Fatal(err)
		}
	}()

	output, err := cmd.CombinedOutput()
	if err != nil {
		return RE, err
	}
	return compareOutput(test.output, string(output)), nil
}

// removeExecutable removes file that creates after compilation
func (ctx *CodeRunnerContext) removeExecutable() {
	err := os.Remove(ctx.executablePath)
	if err != nil {
		log.Fatal(err)
	}
}

// runPartTests runs part of the tests with the specified start and end indexes.
func (ctx *CodeRunnerContext) runPartTests(tests []*Test, start, end, number int) {
	for i := start; i < end; i++ {
		if ctx.failed {
			break
		}
		testResult, err := ctx.runTest(tests[i])
		if err != nil || testResult != OK {
			ctx.failed = true
			ctx.results[number] <- TestingResult{Number: i, Result: testResult, Err: err}
			return
		}
	}
	ctx.results[number] <- TestingResult{Number: -1, Result: OK, Err: nil}
}

// Test tests program on given tests and returns result of testing
func (ctx *CodeRunnerContext) Test(tests []*Test) (TestingResult, error) {
	start := time.Now()
	err := ctx.compileProgram()
	defer ctx.removeExecutable()
	if err != nil {
		return TestingResult{Number: -1, Result: CE}, err
	}
	step := len(tests) / ctx.threads // maybe make constant
	var wg sync.WaitGroup
	for i := 0; i < len(tests); i += step {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx.runPartTests(tests, i, i+step, i/step)
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
