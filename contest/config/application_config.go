// Package config contains structures and functions for contest and application config
package config

import "flag"

// ApplicationConfig represents config for application
// Port is listening port for the web server
// ConfigPath is path to the config for the contest
// LanguagesPath is path to the config for the programming languages
// If Generate is true that application will generate new contest, otherwise that will use existing contest
type ApplicationConfig struct {
	Port          string
	ConfigPath    string
	LanguagesPath string
	Generate      bool
}

// ParseArgs parses arguments from command line
// returns pointer to ApplicationConfig that constructed from command line arguments
func ParseArgs() *ApplicationConfig {
	port := flag.String("port", "8080", "port for the application")
	configPath := flag.String("config", "", "path of the contest config")
	languagesPath := flag.String("languages", "", "path of the languages config")
	generate := flag.Bool(
		"generate",
		false,
		"set if you want to genereate output json file, set to true on the first run",
	)
	flag.Parse()
	return &ApplicationConfig{
		Port:          *port,
		ConfigPath:    *configPath,
		LanguagesPath: *languagesPath,
		Generate:      *generate,
	}
}
