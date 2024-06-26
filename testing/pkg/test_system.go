// Package pkg contains internal logic of the application.
package pkg

import (
	"fmt"
	"log"
	"strings"
	"time"
)

// TestingSystem represents the system that runs tests.
type TestingSystem struct {
	duration           int64
	startTime          string
	testPathForProblem string
	testCount          int
	username           string
	compileFile        string
}

// NewTestingSystem creates new TestingSystem with giving parameters.
func NewTestingSystem(duration int64,
	startTime string,
	testPathForProblem string,
	testCount int,
	username string,
	compileFile string) *TestingSystem {
	return &TestingSystem{
		duration:           duration,
		startTime:          startTime,
		testPathForProblem: testPathForProblem,
		testCount:          testCount,
		username:           username,
		compileFile:        compileFile,
	}
}

// Run represents information that need to execute program on the specified tests
type Run struct {
	fileName string
	language string
	problem  string
	username string
	userId   string
}

// getExecutableName returns name for the executable file and the programming language that was used in the file
func getExecutableName(fileName, language string) string {
	fileName = strings.Replace(fileName, "./", "", 1)
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
func NewRun(fileName, language, problem, username, userId string) *Run {
	return &Run{
		fileName: fileName,
		language: language,
		problem:  problem,
		username: username,
		userId:   userId,
	}
}

// RunTests runs tests and return the result of testing
func (ts *TestingSystem) RunTests(run *Run) (*RunInfo, error) {
	startTime, _ := time.Parse(time.RFC3339, ts.startTime)
	duration := int64(time.Since(startTime).Seconds())
	var returnedErr error = nil
	var result TestingResult
	if duration > ts.duration {
		result = TestingResult{Result: EOC, Number: -1}
	} else {
		executableName := getExecutableName(run.fileName, run.language)
		ctx := NewCodeRunnerContext(run.fileName, executableName, run.language, ts.compileFile)
		result, _ = ctx.Test(ts.testPathForProblem, ts.testCount)
	}
	log.Printf("RESULT: %v", result)
	return NewRunInfo(
		run.userId,
		run.problem,
		result,
		duration,
		run.fileName,
		run.language,
	), returnedErr
}

// RejudgeRun rejudges the specified RunInfo
func (ts *TestingSystem) RejudgeRun(run *RunInfo) error {
	newRunInfo, err := ts.RunTests(NewRun(run.FileName, run.Language, run.Problem, ts.username, run.Id))
	if err != nil {
		return err
	}
	*run = *newRunInfo
	return nil
}
