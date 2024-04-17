package internal

import (
	"fmt"
	"log"
	"strings"
)

// Run represents information that need to execute program on the specified tests
type Run struct {
	directoryWithTests string
	fileName           string
	language           string
}

// getExecutableName returns name for the executable file and the programming language that was used in the file
func getExecutableName(fileName, language string) string {
	arr := strings.Split(fileName, ".")
	if len(arr) != 2 {
		log.Fatal("file name is invalid")
		return ""
	}
	if language == "java" {
		return fmt.Sprintf("%s", arr[0])
	}
	return fmt.Sprintf("%s.%s", arr[0], "out")
}

// NewRun creates Run
func NewRun(fileName, language string) *Run {
	return &Run{
		fileName:           fileName,
		directoryWithTests: "cmd/tests/",
		language:           language,
	}
}

// RunTests runs tests and return the result of testing
func (ts *Run) RunTests() (TestingResult, error) {
	executableName := getExecutableName(ts.fileName, ts.language)
	ctx := NewCodeRunnerContext(ts.fileName, executableName, ts.language)
	return ctx.Test(ts.directoryWithTests, 1000)
}
