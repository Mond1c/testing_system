package internal

import "time"

// ProblemInfo represents number of the task
type ProblemInfo = string

// RunInfo represents information about the running program on the test cases such as the Result and the Time of sending
type RunInfo struct {
	Result TestingResult `json:"result"`
	Time   time.Duration `json:"time"`
}

// NewRunInfo creates pointer of type RunInfo with given RunInfo.Result and RunInfo.Time
func NewRunInfo(result TestingResult, t time.Duration) *RunInfo {
	return &RunInfo{
		Result: result,
		Time:   t,
	}
}

// ContestantInfo represents information about the contestant such as Name, Points, Penalty and information about his Runs.
type ContestantInfo struct {
	Name    string                    `json:"name"`
	Points  int                       `json:"points"`
	Penalty int                       `json:"penalty"`
	Runs    map[ProblemInfo][]RunInfo `json:"runs"`
}

// NewContestantInfo creates pointer of type ContestantInfo
func NewContestantInfo(name string) *ContestantInfo {
	return &ContestantInfo{
		Name:    name,
		Points:  0,
		Penalty: 0,
		Runs:    make(map[ProblemInfo][]RunInfo),
	}
}

// ContestInfo represents information about contest such as Problems names, Contestants information and StartTime of the contest
type ContestInfo struct {
	Problems    []ProblemInfo    `json:"problems"`
	Contestants []ContestantInfo `json:"contestants"`
	StartTime   time.Time        `json:"start_time"`
}

// NewContestInfo creates pointer of type ContestInfo
func NewContestInfo(problems []ProblemInfo, contestants []ContestantInfo, startTime time.Time) *ContestInfo {
	return &ContestInfo{
		Problems:    problems,
		Contestants: contestants,
		StartTime:   startTime,
	}
}
