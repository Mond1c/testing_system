// Package config contains structures and functions for contest and application config
package config

import (
	"encoding/json"
	"errors"
	"os"
)

// Config represents information about the contest.
// TestsInfo represents information about count of the test and path to the directory with tests
// OutputPath represents path to output json file with contest info
type Config struct {
	TestsInfo map[string]struct {
		Path       string `json:"path"`
		TestsCount int    `json:"count"`
	} `json:"tests"`
	OutputPath  string `json:"outputPath"`
	StartTime   string `json:"startTime"`
	Contestants []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"contestants"`
	Credentials map[string]struct {
		Id       string `json:"id"`
		Password string `json:"password"`
		Role     string `json:"role"`
	} `json:"credentials"`
	Problems []string `json:"problems"`
	Duration int64    `json:"duration"`
}

func newConfig() *Config {
	config := Config{}
	config.Duration = 5 * 60 * 60
	config.OutputPath = "output.json"
	return &config
}

func ParseConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := newConfig()
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (c *Config) GetTestPathForProblem(problem string) (string, int, error) {
	value, ok := c.TestsInfo[problem]
	if !ok {
		return "", 0, errors.New("problem doesn't exists")
	}
	return value.Path, value.TestsCount, nil
}

var TestConfig *Config
var TestDir string = "temp"
