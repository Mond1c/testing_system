package main

import (
	"fmt"
	"test_system/internal"
)

func main() {
	//api.InitApi()
	//fmt.Println("Server is starting: 127.0.0.1:8080")
	//err := http.ListenAndServe(":8080", nil)
	//if err != nil {
	//	fmt.Println(err)
	//}
	ctx := internal.NewCodeRunnerContext("cmd/tests/main.cpp")
	result, err := ctx.Test([]*internal.Test{internal.NewTest("1 2", "3"), internal.NewTest("1 3", "4")})
	if err != nil {
		panic(err)
	}
	fmt.Println(result.GetString())
}
