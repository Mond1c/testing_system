package config

import (
	"encoding/json"
	"flag"
	"os"
)

type Config struct {
	Port      string `json:"port"`
	TestDir   string `json:"testDir"`
	TestsInfo map[string]struct {
		Path       string `json:"path"`
		TestsCount int    `json:"count"`
	} `json:"tests"`
	CompileFiles map[string]string `json:"compileFiles"`
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

var MyConfig *Config

type ApplicationConfig = string

func ParseArgs() *ApplicationConfig {
	c := flag.String("config", "config.json", "path to the config file")
	flag.Parse()
	return c
}
