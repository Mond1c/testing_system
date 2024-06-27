package config

import (
	"encoding/json"
	"flag"
	"os"
)

// Config contains configuration for the remote worker
type Config struct {
	Port      string `json:"port"`
	TestDir   string `json:"testDir"`
	TestsInfo map[string]struct {
		Path       string `json:"path"`
		TestsCount int    `json:"count"`
	} `json:"tests"`
	CompileFiles map[string]string `json:"compileFiles"`
}

// newConfig creates new Config
func newConfig() *Config {
	return &Config{}
}

// ParseConfig parses configuration from the specified path
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

// ApplicationConfig represents config for application
type ApplicationConfig = string

// ParseArgs parses arguments from command line
func ParseArgs() *ApplicationConfig {
	c := flag.String("config", "config.json", "path to the config file")
	flag.Parse()
	return c
}
