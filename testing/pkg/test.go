// Package internal contains internal logic of the application.
package pkg

// TestResult type of the enum
type TestResult int

// Enum values
const (
	NONE = TestResult(iota)
	OK
	CE
	RE
	TL
	ML
	WA
	EOC
)

// GetString returns string representation of the test result
func (t *TestResult) GetString() string {
	switch *t {
	case OK:
		return "OK"
	case CE:
		return "Compilation error"
	case RE:
		return "Runtime error"
	case TL:
		return "Time limit"
	case ML:
		return "Memory limit"
	case EOC:
		return "End of the contest"
	case WA:
		return "Wrong answer"
	default:
		return "Unexpected result"
	}
}

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
