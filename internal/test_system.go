package internal

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Mond1c/testing_system/config"
)

// Run represents information that need to execute program on the specified tests
type Run struct {
	fileName string
	language string
	problem  string
	username string
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
func NewRun(fileName, language, problem, username string) *Run {
	return &Run{
		fileName: fileName,
		language: language,
		problem:  problem,
		username: username,
	}
}

// RunTests runs tests and return the result of testing
func (ts *Run) RunTests() (TestingResult, error) {
	startTime, _ := time.Parse(time.RFC3339, config.TestConfig.StartTime)
	duration := int64(time.Since(startTime).Seconds())
	if duration > config.TestConfig.Duration {
		return TestingResult{Result: EOC, Number: -1}, nil
	}
	executableName := getExecutableName(ts.fileName, ts.language)
	ctx := NewCodeRunnerContext(ts.fileName, executableName, ts.language)
	path, count, err := config.TestConfig.GetTestPathForProblem(ts.problem)
	if err != nil {
		return TestingResult{Result: NONE, Number: -1}, err
	}
	result, _ := ctx.Test(path, count)
	log.Printf("RESULT: %v", result)
	AddRun(
		Contest,
		NewRunInfo(
			config.TestConfig.Credentials[ts.username].Id,
			ts.problem,
			result,
			duration,
			ts.fileName,
			ts.language,
		),
	)
	return result, nil
}

func getUsernameById(id string) string {
	for k, v := range config.TestConfig.Credentials {
		if v.Id == id {
			return k
		}
	}
	return ""
}

func RejudgeRun(run *RunInfo) error {
	for i, r := range Contest.Contestants[run.Id].Runs {
		if r.FileName == run.FileName { // i don't really like this, but it will work because I generate uniquie file names
			ts := NewRun(run.FileName, run.Language, run.Problem, getUsernameById(run.Id))
			result, err := ts.RunTests()
			if err != nil {
				return err
			}
			Contest.Contestants[run.Id].Runs[i].Result = result
			return nil
		}
	}
	return errors.New("Can't find run")
}
