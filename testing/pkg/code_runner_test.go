package pkg

import (
	"testing"
)

func TestTestingResult_String(t *testing.T) {
	input := TestingResult{-1, OK}
	expected := "Test with number -1: OK"
	if input.String() != expected {
		t.Errorf("Testing failed, expected %s, got %s", expected, input.String())
	}
}
