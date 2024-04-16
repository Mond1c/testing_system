package internal

import "time"

type ProblemInfo = string

type RunInfo struct {
	Result TestingResult `json:"result"`
	Time   time.Duration `json:"time"`
}

func NewRunInfo(result TestingResult, t time.Duration) *RunInfo {
	return &RunInfo{
		Result: result,
		Time:   t,
	}
}

type ContestantInfo struct {
	Name    string                    `json:"name"`
	Points  int                       `json:"points"`
	Penalty int                       `json:"penalty"`
	Runs    map[ProblemInfo][]RunInfo `json:"runs"`
}

func NewContestantInfo(name string) *ContestantInfo {
	return &ContestantInfo{
		Name:    name,
		Points:  0,
		Penalty: 0,
		Runs:    make(map[ProblemInfo][]RunInfo),
	}
}

type ContestInfo struct {
	Problems    []ProblemInfo    `json:"problems"`
	Contestants []ContestantInfo `json:"contestants"`
	StartTime   time.Time        `json:"start_time"`
}

func NewContestInfo(problems []ProblemInfo, contestants []ContestantInfo, startTime time.Time) *ContestInfo {
	return &ContestInfo{
		Problems:    problems,
		Contestants: contestants,
		StartTime:   startTime,
	}
}
