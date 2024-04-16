package internal

type Test struct {
	input  string
	output string
}

// NewTest creates new test case that can be used in CodeRunnerContext.Test
func NewTest(input, output string) *Test {
	return &Test{input: input, output: output}
}

// TestResult type of the enum
type TestResult int

const (
	OK = TestResult(iota)
	CE
	RE
	TL
	ML
	UB
	WA
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
	case UB:
		return "Something went wrong"
	default:
		return "Unexpected result"
	}
}
