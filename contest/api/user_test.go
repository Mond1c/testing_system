package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Mond1c/testing_system/contest/config"
	"github.com/Mond1c/testing_system/contest/internal"
)

func generateConfigForTest(t *testing.T, startTime time.Time) {
	config.TestConfig = &config.Config{
		TestsInfo: map[string]struct {
			Path       string `json:"path"`
			TestsCount int    `json:"count"`
		}{
			"A": {
				Path:       "../examples/tests",
				TestsCount: 1,
			},
			"B": {
				Path:       "../examples/tests",
				TestsCount: 1,
			},
		},
		OutputPath: "nothing",
		StartTime:  startTime.Format(time.RFC3339),
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
		Duration: 18000,
	}
	contestants := make(map[string]*internal.ContestantInfo)
	for _, contestant := range config.TestConfig.Contestants {
		contestants[contestant.Id] = internal.NewContestantInfo(contestant.Id, contestant.Name)
	}
	startTime, err := time.Parse(time.RFC3339, config.TestConfig.StartTime)
	if err != nil {
		t.Fatal(err)
	}
	internal.Contest = internal.NewContestInfo(config.TestConfig.Problems, contestants, startTime)

	config.LangaugesConfig, err = config.ParseLangauges("../examples/languages.json")
	if err != nil {
		t.Fatal(err)
	}
}

// Fuck this
/*
func TestTest(t *testing.T) {
	generateConfigForTest(t)
	config.TestDir = "."
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	file, err := writer.CreateFormFile("file", "test.cpp")
	if err != nil {
		t.Fatalf("could not create form file: %v", err)
	}
	_, err = io.WriteString(file, `
#include <iostream>
using namespace std;

int main() {
	int a, b;
	cin >> a >> b;
	cout << a + b;
}
`)
	if err != nil {
		t.Fatalf("could not create form file: %v", err)
	}
	_ = writer.WriteField("language", "cpp")
	_ = writer.WriteField("problem", "A")
	_ = writer.WriteField("username", "test1")
	_ = writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/test", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	err = test(rr, req)
	if err != nil {
		t.Fatalf("test failed: %v", err)
	}

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")
	expected := `{"message":"Test with number -1: OK"}`
	assert.JSONEq(t, expected, rr.Body.String(), "Response body mismatch")
}
*/

func TestGetProblems(t *testing.T) {
	startTime := time.Now().UTC()
	generateConfigForTest(t, startTime)
	req := httptest.NewRequest(http.MethodGet, "/api/problems", nil)
	rr := httptest.NewRecorder()
	err := getProblems(rr, req)
	if err != nil {
		t.Errorf("getProblems failed: %v", err)
	}
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")
	expected := `["A","B"]`
	assert.JSONEq(t, expected, rr.Body.String(), "Response body mismatch")
}

func TestGetResults(t *testing.T) {
	startTime := time.Now().UTC()
	generateConfigForTest(t, startTime)
	internal.AddRun(
		internal.Contest,
		internal.NewRunInfo(
			"1",
			"A",
			internal.TestingResult{Result: internal.OK, Number: -1},
			0,
			"test.cpp",
			"cpp",
		),
	)
	req := httptest.NewRequest(http.MethodGet, "/api/results", nil)
	rr := httptest.NewRecorder()
	err := getResults(rr, req)
	if err != nil {
		t.Errorf("getResults failed: %v", err)
	}
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")
	// I don't really like this, but OK
	expected := `
{
"contestants": {
    "1": {
        "additionalPenalty": {},
        "id":"1",
        "name":"Test 1",
        "penalty":0,
        "points":1,
        "results": {
            "A": {
                "fileName":"test.cpp",
                "id":"1",
                "language":"cpp",
                "problem":"A",
                "result":{
                    "number":-1,
                    "result":1
                },
                "time":0
            }
        },
        "runs": [
            {
            "fileName":"test.cpp",
            "id":"1",
            "language":"cpp",
            "problem":"A",
            "result": {
                "number":-1,
                "result":1
            },
            "time":0
            }
        ]
    },
    "2": {
        "additionalPenalty": {},
        "id":"2",
        "name":"Test 2",
        "penalty": 0,
        "points": 0,
        "results": {},
        "runs": []
    }},
"problems": ["A", "B"],
"startTime":"timeHere"
}`
	expected = strings.Replace(expected, "timeHere", startTime.Format(time.RFC3339), 1)
	assert.JSONEq(t, expected, rr.Body.String(), "Response body mismatch")
}

func TestGetLanguages(t *testing.T) {
	startTime := time.Now().UTC()
	generateConfigForTest(t, startTime)
	req := httptest.NewRequest(http.MethodGet, "/api/languages", nil)
	rr := httptest.NewRecorder()
	err := getLanguages(rr, req)
	if err != nil {
		t.Errorf("getLanguages failed: %v", err)
	}
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")
	expected := `
[
    {
        "name": "C++ 20",
        "value": "cpp"
    },
    {
        "name": "Go 1.21",
        "value": "go"
    },
    {
        "name": "Java 21",
        "value": "java"
    }
]
`
	assert.JSONEq(t, expected, rr.Body.String(), "Response body mismatch")
}

func TestGetContestStartTime(t *testing.T) {
	startTime := time.Now().UTC()
	generateConfigForTest(t, startTime)
	req := httptest.NewRequest(http.MethodGet, "/api/startTime", nil)
	rr := httptest.NewRecorder()
	err := getContestStartTime(rr, req)

	if err != nil {
		t.Errorf("getContestStartTime failed: %v", err)
	}

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")

	expectedTime, _ := time.Parse(time.RFC3339, startTime.Format(time.RFC3339))
	expected := StartTimeResponse{
		StartTime: expectedTime.UnixMilli(),
		Duration:  18000,
	}
	var actual StartTimeResponse
	err = json.Unmarshal(rr.Body.Bytes(), &actual)
	assert.NoError(t, err, "Failed to unmarshal response")
	assert.Equal(t, expected, actual, "Response body mismatch")
}
