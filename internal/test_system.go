package internal

import (
	"fmt"
	"log"
	"strings"
)

type TestSystem struct {
	tests    []*Test
	fileName string
}

func generateTests() []*Test {
	tests := make([]*Test, 1000)
	for i := 0; i < 1000; i++ {
		tests[i] = NewTest(fmt.Sprintf("%d %d", i, i+1), fmt.Sprintf("%d", i+i+1))
	}
	return tests
}

func getExecutableName(fileName string) string {
	arr := strings.Split(fileName, ".")
	if len(arr) < 2 {
		log.Fatal("file name is invalid")
	}
	return fmt.Sprintf("%s.%s", arr[0], ".out")
}

func NewTestSystem(fileName string) *TestSystem {
	return &TestSystem{
		fileName: fileName,
		tests:    generateTests(),
	}
}

func (ts *TestSystem) Run() (TestingResult, error) {
	ctx := NewCodeRunnerContext(ts.fileName, getExecutableName(ts.fileName))
	return ctx.Test(ts.tests)
}
