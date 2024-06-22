package internal

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"test_system/config"
)

// ProblemInfo represents number of the task
type ProblemInfo = string

// RunInfo represents information about the run such as Id, Problem, Result, Time and FileName
type RunInfo struct {
	Id       string        `json:"id"`
	Problem  ProblemInfo   `json:"problem"`
	Result   TestingResult `json:"result"`
	Time     int64         `json:"time"`
	FileName string        `json:"fileName"`
	Language string        `json:"language"`
}

// NewRunInfo creates pointer of type RunInfo with given RunInfo.Result and RunInfo.Time
func NewRunInfo(id, problem string, result TestingResult, t int64, fileName, language string) *RunInfo {
	return &RunInfo{
		Id:       id,
		Problem:  problem,
		Result:   result,
		Time:     t,
		FileName: fileName,
		Language: language,
	}
}

// ContestantInfo represents information about the contestant such as Name, Points, Penalty and information about his Runs.
type ContestantInfo struct {
	Id                string                  `json:"id"`
	Name              string                  `json:"name"`
	Points            int                     `json:"points"`
	Penalty           int64                   `json:"penalty"`
	Runs              []RunInfo               `json:"runs"`
	Results           map[ProblemInfo]RunInfo `json:"results"`
	AdditionalPenalty map[ProblemInfo]int64   `json:"additionalPenalty"`
	mu                sync.Mutex
}

// NewContestantInfo creates pointer of type ContestantInfo
func NewContestantInfo(id, name string) *ContestantInfo {
	return &ContestantInfo{
		Id:                id,
		Name:              name,
		Points:            0,
		Penalty:           0,
		Runs:              make([]RunInfo, 0),
		Results:           make(map[ProblemInfo]RunInfo),
		AdditionalPenalty: make(map[ProblemInfo]int64),
	}
}

// ContestInfo represents information about contest such as Problems names, Contestants information and StartTime of the contest
type ContestInfo struct {
	Problems    []ProblemInfo              `json:"problems"`
	Contestants map[string]*ContestantInfo `json:"contestants"`
	StartTime   time.Time                  `json:"startTime"`
}

// NewContestInfo creates pointer of type ContestInfo
func NewContestInfo(
	problems []ProblemInfo,
	contestants map[string]*ContestantInfo,
	startTime time.Time,
) *ContestInfo {
	return &ContestInfo{
		Problems:    problems,
		Contestants: contestants,
		StartTime:   startTime,
	}
}

const timeStepForContestUpdateMs = 10000

// AddRun adds new run for the current contest
func AddRun(contest *ContestInfo, run *RunInfo) {
	prevResult := contest.Contestants[run.Id].Results[run.Problem]
	contest.Contestants[run.Id].mu.Lock()

	if prevResult.Result.Result != OK && run.Result.Result == OK {
		contest.Contestants[run.Id].Results[run.Problem] = *run
		contest.Contestants[run.Id].Points += 1
		contest.Contestants[run.Id].Penalty += run.Time + contest.Contestants[run.Id].AdditionalPenalty[run.Problem]
	} else if prevResult.Result.Result != OK && run.Result.Result != OK {
		log.Print(contest.Contestants[run.Id].AdditionalPenalty)
		// TODO: 20 is Constant?
		contest.Contestants[run.Id].AdditionalPenalty[run.Problem] += 20
	}

	contest.Contestants[run.Id].Runs = append(contest.Contestants[run.Id].Runs, *run)
	contest.Contestants[run.Id].mu.Unlock()
}

// GenerateContestInfo generates default contest info json file with the specified config
func GenerateContestInfo(config *config.Config) error {
	startTime, err := time.Parse(time.RFC3339, config.StartTime)
	if err != nil {
		log.Printf("Can't generate contest info becase start time is invalid: %v", err)
		return err
	}
	contestants := make(map[string]*ContestantInfo)
	for _, contestant := range config.Contestants {
		contestants[contestant.Id] = NewContestantInfo(contestant.Id, contestant.Name)
	}

	contest := NewContestInfo(config.Problems, contestants, startTime)
	data, err := json.Marshal(contest)
	if err != nil {
		log.Printf("Can't encode contest info to json: %v", err)
		return err
	}
	err = os.WriteFile(config.OutputPath, data, 0644)
	if err != nil {
		log.Printf("Can't write contest info to file: %v", err)
	}
	return err
}

// UpdateContestInfo updates info about the current contest and writes it to the specified json file
func UpdateContestInfo(config *config.Config, contest **ContestInfo) {
	if *contest == nil {
		data, err := os.ReadFile(config.OutputPath)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(data, contest)
		if err != nil {
			log.Fatal(err)
		}
	}
	for {
		log.Println("Starting update contest info!")
		data, err := json.Marshal(*contest)
		if err == nil {
			err = os.WriteFile(config.OutputPath, data, 0644)
		}
		if err != nil {
			log.Printf("Failed to update contest info: %v", err)
		}
		log.Println("Ending update contest info!")
		time.Sleep(time.Millisecond * timeStepForContestUpdateMs)
	}
}

var Contest *ContestInfo
