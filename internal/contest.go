package internal

import (
	"encoding/json"
	"log"
	"os"
	"test_system/config"
	"time"
)

// ProblemInfo represents number of the task
type ProblemInfo = string

// RunInfo represents information about the running program on the test cases such as the Result and the Time of sending
type RunInfo struct {
	Id      string        `json:"id"`
	Problem ProblemInfo   `json:"problem"`
	Result  TestingResult `json:"result"`
	Time    time.Duration `json:"time"`
}

// NewRunInfo creates pointer of type RunInfo with given RunInfo.Result and RunInfo.Time
func NewRunInfo(id string, result TestingResult, t time.Duration) *RunInfo {
	return &RunInfo{
		Id:     id,
		Result: result,
		Time:   t,
	}
}

// ContestantInfo represents information about the contestant such as Name, Points, Penalty and information about his Runs.
type ContestantInfo struct {
	Id      string                  `json:"id"`
	Name    string                  `json:"name"`
	Points  int                     `json:"points"`
	Penalty int                     `json:"penalty"`
	Runs    []RunInfo               `json:"runs"`
	Results map[ProblemInfo]RunInfo `json:"results"`
}

// NewContestantInfo creates pointer of type ContestantInfo
func NewContestantInfo(id, name string) *ContestantInfo {
	return &ContestantInfo{
		Id:      id,
		Name:    name,
		Points:  0,
		Penalty: 0,
		Runs:    make([]RunInfo, 0),
	}
}

// ContestInfo represents information about contest such as Problems names, Contestants information and StartTime of the contest
type ContestInfo struct {
	Problems    []ProblemInfo             `json:"problems"`
	Contestants map[string]ContestantInfo `json:"contestants"`
	StartTime   time.Time                 `json:"start_time"`
}

// NewContestInfo creates pointer of type ContestInfo
func NewContestInfo(problems []ProblemInfo, contestants map[string]ContestantInfo, startTime time.Time) *ContestInfo {
	return &ContestInfo{
		Problems:    problems,
		Contestants: contestants,
		StartTime:   startTime,
	}
}

const timeStepForContestUpdateMs = 10000

var runs []*RunInfo

func AddRun(run *RunInfo) {
	runs = append(runs, run)
}

func UpdateContestInfo() {
	for {
		data, err := os.ReadFile(config.TestConfig.OutputPath)
		if err != nil {
			log.Fatalf("Failed to update contest info: %v", err)
			return
		}
		contest := &ContestInfo{}
		err = json.Unmarshal(data, &contest)
		if err != nil {
			log.Fatalf("Failed to update contest info: %v", err)
			return
		}
		for _, run := range runs {
			prevResult := contest.Contestants[run.Id].Results[run.Problem]
			if prevResult.Result.Result != OK && run.Result.Result == OK {
				contest.Contestants[run.Id].Results[run.Problem] = *run
			}
		}
		data, err = json.Marshal(*contest)
		if err != nil {
			log.Fatalf("Failed to update contest info: %v", err)
			return
		}
		err = os.WriteFile(config.TestConfig.OutputPath, data, 0644)
		if err != nil {
			log.Fatalf("Failed to update contest info: %v", err)
			return
		}
		time.Sleep(time.Millisecond * timeStepForContestUpdateMs)
	}
}
