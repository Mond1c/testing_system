// Package internal contains internal logic of the application.
package internal

import (
	"encoding/json"
	"github.com/Mond1c/testing_system/contest/config"
	"github.com/Mond1c/testing_system/testing/pkg"
	"log"
	"os"
	"sync"
	"time"
)

// ProblemInfo represents number of the task
type ProblemInfo = string

// ContestantInfo represents information about the contestant such as Name, Points, Penalty and information about his Runs.
type ContestantInfo struct {
	Id                string                      `json:"id"`
	Name              string                      `json:"name"`
	Points            int                         `json:"points"`
	Penalty           int64                       `json:"penalty"`
	Runs              []pkg.RunInfo               `json:"runs"`
	Results           map[ProblemInfo]pkg.RunInfo `json:"results"`
	AdditionalPenalty map[ProblemInfo]int64       `json:"additionalPenalty"`
	mu                sync.Mutex
}

// NewContestantInfo creates pointer of type ContestantInfo
func NewContestantInfo(id, name string) *ContestantInfo {
	return &ContestantInfo{
		Id:                id,
		Name:              name,
		Points:            0,
		Penalty:           0,
		Runs:              make([]pkg.RunInfo, 0),
		Results:           make(map[ProblemInfo]pkg.RunInfo),
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
const penaltyForWrongAnswer = 20

// AddRun adds new run for the current contest
func AddRun(contest *ContestInfo, run *pkg.RunInfo) {
	prevResult := contest.Contestants[run.Id].Results[run.Problem]
	contest.Contestants[run.Id].mu.Lock()

	if prevResult.Result.Result != pkg.OK && run.Result.Result == pkg.OK {
		contest.Contestants[run.Id].Results[run.Problem] = *run
		contest.Contestants[run.Id].Points += 1
		contest.Contestants[run.Id].Penalty += run.Time + contest.Contestants[run.Id].AdditionalPenalty[run.Problem]
	} else if prevResult.Result.Result != pkg.OK && run.Result.Result != pkg.OK {
		log.Print(contest.Contestants[run.Id].AdditionalPenalty)
		contest.Contestants[run.Id].AdditionalPenalty[run.Problem] += penaltyForWrongAnswer
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
