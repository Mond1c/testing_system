package config

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	TestsInfo map[string]struct {
		Path       string `json:"path"`
		TestsCount int    `json:"count"`
	} `json:"tests"`
	OutputPath string `json:"outputPath"`
	StartTime  string `json:"startTime"`
	Contestans []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"contestants"`
	Problems []string `json:"problems"`
}

func newConfig() *Config {
	return &Config{}
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
