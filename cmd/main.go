package main

import (
	"log"
	"net/http"
	"test_system/api"
)

func main() {
	api.InitApi()
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
		return
	}
}
