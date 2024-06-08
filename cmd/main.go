package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"test_system/api"
	"test_system/config"
	"test_system/internal"
)

func CheckIfFileExists(path string) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		log.Fatal(err)
	}
}

func parseConfig[T any](f func(path string) (T, error), path string) T {
	c, err := f(path)
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func main() {
	applicationConfig := config.ParseArgs()

	CheckIfFileExists(applicationConfig.ConfigPath)
	CheckIfFileExists(applicationConfig.LanguagesPath)

	config.TestConfig = parseConfig(config.ParseConfig, applicationConfig.ConfigPath)
	config.LangaugesConfig = parseConfig(config.ParseLangauges, applicationConfig.LanguagesPath)

	log.Printf("Names: %v", config.LangaugesConfig.GetLanguages())

	if _, err := os.Stat(config.TestDir); !errors.Is(err, os.ErrNotExist) {
		log.Printf("Directory %s already exists", config.TestDir)
	} else {
		err := os.Mkdir(config.TestDir, 0750)
		if err != nil {
			log.Fatal(err)
		}
	}

	if applicationConfig.Generate {
		err := internal.GenerateContestInfo()
		if err != nil {
			return
		}
	}
	fs := http.FileServer(http.Dir("./frontend/build"))
	http.Handle("/static/", fs)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || filepath.Ext(r.URL.Path) == "" {
			http.ServeFile(w, r, "./frontend/build/index.html")
		} else {
			fs.ServeHTTP(w, r)
		}
	})

	api.InitUserApi()
	api.InitAdminAPI()
	go internal.UpdateContestInfo()
	err := http.ListenAndServe(":"+applicationConfig.Port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
