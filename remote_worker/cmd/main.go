package main

import (
	"errors"
	"github.com/Mond1c/testing_system/remote_worker/api"
	"github.com/Mond1c/testing_system/remote_worker/config"
	"log"
	"net/http"
	"os"
)

func main() {
	ac := config.ParseArgs()
	var err error
	config.MyConfig, err = config.ParseConfig(*ac)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stat(config.MyConfig.TestDir); !errors.Is(err, os.ErrNotExist) {
		log.Printf("Directory %s already exists", config.MyConfig.TestDir)
	} else {
		err := os.Mkdir(config.MyConfig.TestDir, 0750)
		if err != nil {
			log.Fatal(err)
		}
	}
	api.InitAPI()
	log.Print(http.ListenAndServe(":"+config.MyConfig.Port, nil))
}
