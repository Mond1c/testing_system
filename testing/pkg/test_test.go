package pkg

import "testing"

func TestGetStringMethodForTestResult(t *testing.T) {
	tests := []struct {
		input    TestResult
		expected string
	}{
		{NONE, "Unexpected result"},
		{OK, "OK"},
		{CE, "Compilation error"},
		{RE, "Runtime error"},
		{TL, "Time limit"},
		{ML, "Memory limit"},
		{WA, "Wrong answer"},
		{EOC, "End of the contest"},
	}

	for _, tt := range tests {
		if tt.expected != tt.input.GetString() {
			t.Errorf("GetStringMethodForTestResult.String() expected %v, got %v", tt.expected, tt.input.GetString())
		}
	}
}
