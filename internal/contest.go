package internal

import (
	"encoding/json"
	"log"
	"os"
	"sync"
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
	Id               string                  `json:"id"`
	Name             string                  `json:"name"`
	Points           int                     `json:"points"`
	Penalty          int64                   `json:"penalty"`
	Runs             []RunInfo               `json:"runs"`
	Results          map[ProblemInfo]RunInfo `json:"results"`
	AdditinalPenalty map[ProblemInfo]int64   `json:"additionalPenalty"`
	mu               sync.Mutex
}

// NewContestantInfo creates pointer of type ContestantInfo
func NewContestantInfo(id, name string) *ContestantInfo {
	return &ContestantInfo{
		Id:               id,
		Name:             name,
		Points:           0,
		Penalty:          0,
		Runs:             make([]RunInfo, 0),
		Results:          make(map[ProblemInfo]RunInfo),
		AdditinalPenalty: make(map[ProblemInfo]int64),
	}
}

// ContestInfo represents information about contest such as Problems names, Contestants information and StartTime of the contest
type ContestInfo struct {
	Problems    []ProblemInfo              `json:"problems"`
	Contestants map[string]*ContestantInfo `json:"contestants"`
	StartTime   time.Time                  `json:"start_time"`
}

// NewContestInfo creates pointer of type ContestInfo
func NewContestInfo(problems []ProblemInfo, contestants map[string]*ContestantInfo, startTime time.Time) *ContestInfo {
	return &ContestInfo{
		Problems:    problems,
		Contestants: contestants,
		StartTime:   startTime,
	}
}

const timeStepForContestUpdateMs = 10000

// AddRun adds new run for the current contest
func AddRun(run *RunInfo) {
	prevResult := Contest.Contestants[run.Id].Results[run.Problem]
	Contest.Contestants[run.Id].mu.Lock()

	if prevResult.Result.Result != OK && run.Result.Result == OK {
		Contest.Contestants[run.Id].Results[run.Problem] = *run
		Contest.Contestants[run.Id].Points += 1
		Contest.Contestants[run.Id].Penalty += run.Time + Contest.Contestants[run.Id].AdditinalPenalty[run.Problem]
	} else if prevResult.Result.Result != OK && run.Result.Result != OK {
		log.Print(Contest.Contestants[run.Id].AdditinalPenalty)
		Contest.Contestants[run.Id].AdditinalPenalty[run.Problem] += 20
	}

	Contest.Contestants[run.Id].Runs = append(Contest.Contestants[run.Id].Runs, *run)
	Contest.Contestants[run.Id].mu.Unlock()
}

// GenerateContestInfo generates default contest info json file with the specified config
func GenerateContestInfo() error {
	startTime, err := time.Parse(time.RFC3339, config.TestConfig.StartTime)
	if err != nil {
		log.Printf("Can't generate contest info becase invalid start time: %v", err)
		return err
	}
	contestants := make(map[string]*ContestantInfo)
	for _, contestant := range config.TestConfig.Contestans {
		contestants[contestant.Id] = NewContestantInfo(contestant.Id, contestant.Name)
	}
	contest := NewContestInfo(config.TestConfig.Problems, contestants, startTime)
	log.Printf("%v\n", contest)
	data, err := json.Marshal(contest)
	if err != nil {
		log.Printf("Can't encode contest info to json: %v", err)
		return err
	}
	err = os.WriteFile(config.TestConfig.OutputPath, data, 0644)
	if err != nil {
		log.Printf("Can't write contest info to file: %v", err)
	}
	return err
}

// TODO: try to decrease memory allocation (maybe use jsonparse)
// UpdateContestInfo upates info about the current contest and writes it to the specified json file
func UpdateContestInfo() {
	if Contest == nil {
		data, err := os.ReadFile(config.TestConfig.OutputPath)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(data, &Contest)
		if err != nil {
			log.Fatal(err)
		}
	}
	for {
		log.Println("Starting update contest info!")
		data, err := json.Marshal(*Contest)
		if err != nil {
			log.Printf("Failed to update contest info: %v", err)
			return
		}
		err = os.WriteFile(config.TestConfig.OutputPath, data, 0644)
		if err != nil {
			log.Printf("Failed to update contest info: %v", err)
			return
		}
		runs = nil
		log.Println("Ending update contest info!")
		time.Sleep(time.Millisecond * timeStepForContestUpdateMs)
	}
}

var Contest *ContestInfo
