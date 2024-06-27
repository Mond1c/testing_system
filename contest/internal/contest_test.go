package internal

import (
	"encoding/json"
	"github.com/Mond1c/testing_system/contest/config"
	"github.com/Mond1c/testing_system/testing/pkg"
	"os"
	"strconv"
	"testing"
	"time"
)

func runInfoEquals(lhs, rhs *pkg.RunInfo) bool {
	return lhs.Id == rhs.Id && lhs.Problem == rhs.Problem && lhs.Time == rhs.Time && lhs.FileName == rhs.FileName &&
		lhs.Result.Result == rhs.Result.Result && lhs.Result.Number == rhs.Result.Number && lhs.Language == rhs.Language
}

func problemEquals(lhs []ProblemInfo, rhs []string) bool {
	if len(lhs) != len(rhs) {
		return false
	}
	for i := range lhs {
		if lhs[i] != rhs[i] {
			return false
		}
	}
	return true
}

func TestAddRun(t *testing.T) {
	problems := []ProblemInfo{"A", "B", "C"}
	contestants := map[string]*ContestantInfo{
		"1": NewContestantInfo("1", "Test 1"),
		"2": NewContestantInfo("2", "Test 2"),
	}
	contest := NewContestInfo(problems, contestants, time.Now())
	runInfo := pkg.NewRunInfo("1", "A", pkg.TestingResult{-1, pkg.OK}, 10, "test.cpp", "cpp")
	AddRun(contest, runInfo)

	if len(contest.Contestants["1"].Runs) != 1 {
		t.Error("Expected 1 run")
	}
	if !runInfoEquals(&contest.Contestants["1"].Runs[0], runInfo) {
		t.Error("Expected run result")
	}
	if contest.Contestants["1"].Points != 1 {
		t.Error("Expected 1 point")
	}
	if contest.Contestants["1"].Penalty != 10 {
		t.Error("Expected 10 penalty")
	}
}

func TestAddRunWithExtraPenalty(t *testing.T) {
	problems := []ProblemInfo{"A", "B", "C"}
	contestants := map[string]*ContestantInfo{
		"1": NewContestantInfo("1", "Test 1"),
		"2": NewContestantInfo("2", "Test 2"),
	}
	contest := NewContestInfo(problems, contestants, time.Now())
	badRunInfo := pkg.NewRunInfo("1", "A", pkg.TestingResult{0, pkg.WA}, 10, "test.cpp", "cpp")
	AddRun(contest, badRunInfo)
	runInfo := pkg.NewRunInfo("1", "A", pkg.TestingResult{-1, pkg.OK}, 20, "test.cpp", "cpp")
	AddRun(contest, runInfo)

	if len(contest.Contestants["1"].Runs) != 2 {
		t.Error("Expected 1 run")
	}
	if !runInfoEquals(&contest.Contestants["1"].Runs[0], badRunInfo) {
		t.Error("Expected run result")
	}
	if !runInfoEquals(&contest.Contestants["1"].Runs[1], runInfo) {
		t.Error("Expected run result")
	}
	if contest.Contestants["1"].Points != 1 {
		t.Error("Expected 1 point")
	}
	if contest.Contestants["1"].Penalty != 40 {
		t.Error("Expected 10 penalty")
	}
}

func TestGenerateContestInfo(t *testing.T) {
	file, err := os.CreateTemp(".", "*_test.json")
	if err != nil {
		t.Error(err)
	}
	file.Close()
	defer removeFile(file.Name())

	testConfig := config.Config{
		TestsInfo: map[string]struct {
			Path       string `json:"path"`
			TestsCount int    `json:"count"`
		}{
			"A": {
				Path:       "test1",
				TestsCount: 1,
			},
			"B": {
				Path:       "test2",
				TestsCount: 2,
			},
		},
		OutputPath: file.Name(),
		StartTime:  "2024-05-30T18:00:00Z",
		Contestants: []struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		}{
			{
				Id:   "1",
				Name: "Test 1",
			},
			{
				Id:   "2",
				Name: "Test 2",
			},
		},
		Credentials: map[string]struct {
			Id       string `json:"id"`
			Password string `json:"password"`
			Role     string `json:"role"`
		}{
			"test1": {
				Id:       "1",
				Password: "test1",
			},
			"test2": {
				Id:       "2",
				Password: "test2",
			},
		},
		Problems: []string{"A", "B"},
	}

	err = GenerateContestInfo(&testConfig)
	if err != nil {
		t.Error(err)
	}

	var contest *ContestInfo
	data, err := os.ReadFile(file.Name())
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal(data, &contest)
	if err != nil {
		t.Error(err)
	}
	if !problemEquals(contest.Problems, testConfig.Problems) {
		t.Error("Expected problems")
	}

	if len(contest.Contestants) != len(testConfig.Credentials) {
		t.Error("Expected number of contestants")
	}
	for i, c := range contest.Contestants {
		index, err := strconv.Atoi(i)
		if err != nil {
			t.Error(err)
		}
		contestant := testConfig.Contestants[index-1]
		if c.Id != contestant.Id || c.Name != contestant.Name {
			t.Error("Expected contestant")
		}
	}
	st, err := time.Parse(time.RFC3339, testConfig.StartTime)
	if err != nil {
		t.Error(err)
	}
	if contest.StartTime != st {
		t.Error("Expected start time")
	}
}

func TestGenerateContestInfoWithInvalidTime(t *testing.T) {
	file, err := os.CreateTemp(".", "*_test.json")
	if err != nil {
		t.Error(err)
	}
	file.Close()
	defer removeFile(file.Name())

	testConfig := config.Config{
		TestsInfo: map[string]struct {
			Path       string `json:"path"`
			TestsCount int    `json:"count"`
		}{
			"A": {
				Path:       "test1",
				TestsCount: 1,
			},
			"B": {
				Path:       "test2",
				TestsCount: 2,
			},
		},
		OutputPath: file.Name(),
		StartTime:  "123",
		Contestants: []struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		}{
			{
				Id:   "1",
				Name: "Test 1",
			},
			{
				Id:   "2",
				Name: "Test 2",
			},
		},
		Credentials: map[string]struct {
			Id       string `json:"id"`
			Password string `json:"password"`
			Role     string `json:"role"`
		}{
			"test1": {
				Id:       "1",
				Password: "test1",
			},
			"test2": {
				Id:       "2",
				Password: "test2",
			},
		},
		Problems: []string{"A", "B"},
	}

	err = GenerateContestInfo(&testConfig)
	if err == nil {
		t.Error("Expected error")
	}
}
