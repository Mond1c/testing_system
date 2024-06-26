// Package internal contains internal logic of the application.
package internal

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
