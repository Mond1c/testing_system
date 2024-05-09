package config

import "flag"

type ApplicationConfig struct {
	Port          string
	ConfigPath    string
	LanguagesPath string
	Generate      bool
}

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
