package internal

// TestResult type of the enum
type TestResult int

const (
	NONE = TestResult(iota)
	OK
	CE
	RE
	TL
	ML
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
	default:
		return "Unexpected result"
	}
}
