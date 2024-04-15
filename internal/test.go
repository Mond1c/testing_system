package internal

type Test struct {
	input  string
	output string
}

func NewTest(input, output string) *Test {
	return &Test{input: input, output: output}
}

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

// TODO: Maybe return error

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
	}
	return "Unexpected result"
}
