package internal

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type CodeRunnerContext struct {
	filePath       string
	language       string
	executablePath string
}

func NewCodeRunnerContext(filePath string) *CodeRunnerContext {
	return &CodeRunnerContext{
		filePath:       filePath,
		language:       "c++",
		executablePath: "test.out",
	}
}

type TestingResult struct {
	number int
	result TestResult
}

func (t *TestingResult) GetString() string {
	return fmt.Sprintf("Test with number %d: %s", t.number, t.result.GetString())
}

// compileProgram compiles source code to executable file using giving CodeRunnerContext.
// Using specific compiler based on given language.
func (ctx *CodeRunnerContext) compileProgram() error {
	if _, err := os.Stat(ctx.filePath); os.IsNotExist(err) {
		return err
	}
	_, err := exec.Command("g++", "-std=c++20", "-o", ctx.executablePath, ctx.filePath).Output()
	if err != nil {
		return err
	}
	_, err = exec.Command("chmod", "+x", ctx.executablePath).Output()
	return err
}

func compareOutput(original, output string) TestResult {
	output = strings.Trim(output, "\n")
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
		defer input.Close()
		io.WriteString(input, test.input)
	}()

	output, err := cmd.CombinedOutput()
	if err != nil {
		return RE, err
	}
	return compareOutput(test.output, string(output)), nil
}

// TODO: rewrite to gorutines

// Test tests program on given tests and returns result of testing
func (ctx *CodeRunnerContext) Test(tests []*Test) (TestingResult, error) {
	err := ctx.compileProgram()
	if err != nil {
		return TestingResult{number: -1, result: CE}, err
	}
	for i := 0; i < len(tests); i++ {
		testResult, err := ctx.runTest(tests[i])
		if err != nil {
			return TestingResult{number: i, result: testResult}, err
		}
		if testResult != OK {
			return TestingResult{number: i, result: testResult}, err
		}
	}
	return TestingResult{number: -1, result: OK}, nil
}
