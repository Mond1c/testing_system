package main

import (
	"fmt"
	"net/http"
	"test_system/api"
)

func main() {
	api.InitApi()
	fmt.Println("Server is starting: 127.0.0.1:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
