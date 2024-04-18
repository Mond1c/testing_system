package internal

import (
	"fmt"
	"log"
	"strings"
	"test_system/config"
	"time"
)

// Run represents information that need to execute program on the specified tests
type Run struct {
	directoryWithTests string
	fileName           string
	language           string
	problem            string
	username           string
}

// getExecutableName returns name for the executable file and the programming language that was used in the file
func getExecutableName(fileName, language string) string {
	arr := strings.Split(fileName, ".")
	if len(arr) != 2 {
		log.Println("file name is invalid")
		return ""
	}
	if language == "java" {
		return arr[0]
	}
	return fmt.Sprintf("%s.%s", arr[0], "out")
}

// NewRun creates Run
func NewRun(fileName, language, problem, username string) *Run {
	return &Run{
		fileName:           fileName,
		directoryWithTests: "cmd/tests/",
		language:           language,
		problem:            problem,
		username:           username,
	}
}

// RunTests runs tests and return the result of testing
func (ts *Run) RunTests() (TestingResult, error) {
	sends, _ := time.Parse(time.RFC3339, config.TestConfig.StartTime)
	duration := int64(time.Since(sends).Minutes())
	executableName := getExecutableName(ts.fileName, ts.language)
	ctx := NewCodeRunnerContext(ts.fileName, executableName, ts.language)
	path, count, err := config.TestConfig.GetTestPathForProblem(ts.problem)
	if err != nil {
		log.Println(err)
		return TestingResult{}, err
	}
	result, _ := ctx.Test(path, count)
	log.Printf("RESULT: %v", result)
	AddRun(
		NewRunInfo(LoginContestantId[ts.username], ts.problem, result, duration),
	)
	return result, nil
}
