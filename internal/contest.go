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
	Time    int64         `json:"time"`
}

// NewRunInfo creates pointer of type RunInfo with given RunInfo.Result and RunInfo.Time
func NewRunInfo(id, problem string, result TestingResult, t int64) *RunInfo {
	return &RunInfo{
		Id:      id,
		Problem: problem,
		Result:  result,
		Time:    t,
	}
}

// ContestantInfo represents information about the contestant such as Name, Points, Penalty and information about his Runs.
type ContestantInfo struct {
	Id      string                  `json:"id"`
	Name    string                  `json:"name"`
	Points  int                     `json:"points"`
	Penalty int64                   `json:"penalty"`
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

// AddRun adds new run for the current contest
func AddRun(run *RunInfo) {
	runs = append(runs, run)
}

// GenerateContestInfo generates default contest info json file with the specified config
func GenerateContestInfo() error {
	startTime, err := time.Parse(time.RFC3339, config.TestConfig.StartTime)
	if err != nil {
		log.Fatalf("Can't generate contest info becase invalid start time: %v", err)
		return err
	}
	contestants := make(map[string]ContestantInfo)
	for _, contestant := range config.TestConfig.Contestans {
		contestants[contestant.Id] = *NewContestantInfo(contestant.Id, contestant.Name)
	}
	contest := NewContestInfo(config.TestConfig.Problems, contestants, startTime)
	log.Printf("%v\n", contest)
	data, err := json.Marshal(contest)
	if err != nil {
		log.Fatalf("Can't encode contest info to json: %v", err)
		return err
	}
	err = os.WriteFile(config.TestConfig.OutputPath, data, 0644)
	if err != nil {
		log.Fatalf("Can't write contest info to file: %v", err)
	}
	return err
}

// UpdateContestInfo upates info about the current contest and writes it to the specified json file
func UpdateContestInfo() {
	for {
		log.Println("Starting update contest info!")
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
			contestant := contest.Contestants[run.Id]
			if prevResult.Result.Result != OK && run.Result.Result == OK {
				contestantResults := contestant.Results
				if contestantResults == nil {
					contestantResults = make(map[ProblemInfo]RunInfo)
				}
				contestantResults[run.Problem] = *run
				contestant.Points += 1
				// TODO: Add penalty
				contestant.Penalty += run.Time
				contestant.Results = contestantResults
			} else if prevResult.Result.Result != OK && run.Result.Result != OK {
				contestant.Penalty += 20
			}
			contestantRuns := contestant.Runs
			if contestantRuns == nil {
				contestantRuns = make([]RunInfo, 0)
			}
			contestantRuns = append(contestantRuns, *run)
			contestant.Runs = contestantRuns
			contest.Contestants[run.Id] = contestant
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
		runs = nil
		log.Println("Ending update contest info!")
		time.Sleep(time.Millisecond * timeStepForContestUpdateMs)
	}
}
