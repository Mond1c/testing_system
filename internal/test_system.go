package internal

import (
	"fmt"
	"log"
	"strings"
)

// Run represents information that need to execute program on the specified tests
type Run struct {
	tests    []*Test
	fileName string
}

// generateTests test function
// TODO: need to delete in the future
func generateTests() []*Test {
	tests := make([]*Test, 1000)
	for i := 0; i < 1000; i++ {
		tests[i] = NewTest(fmt.Sprintf("%d %d", i, i+1), fmt.Sprintf("%d", i+i+1))
	}
	return tests
}

// getExecutableNameAndLanguage returns name for the executable file and the programming language that was used in the file
func getExecutableNameAndLanguage(fileName string) (string, string) {
	arr := strings.Split(fileName, ".")
	if len(arr) != 2 {
		log.Fatal("file name is invalid")
		return "", ""
	}
	return fmt.Sprintf("%s.%s", arr[0], "out"), arr[1]
}

// NewRun creates Run
func NewRun(fileName string) *Run {
	return &Run{
		fileName: fileName,
		tests:    generateTests(),
	}
}

// RunTests runs tests and return the result of testing
func (ts *Run) RunTests() (TestingResult, error) {
	executableName, executableLang := getExecutableNameAndLanguage(ts.fileName)
	ctx := NewCodeRunnerContext(ts.fileName, executableName, executableLang)
	return ctx.Test(ts.tests)
}
