package main

import (
	"fmt"
	"log"
	"test_system/internal"
)

func main() {
	ctx := internal.NewCodeRunnerContext("cmd/tests/main.cpp")
	result, err := ctx.Test([]*internal.Test{internal.NewTest("1 2", "3"), internal.NewTest("1 3", "4")})
	if err != nil {
		log.Fatal(result.GetString())
	}
	fmt.Println(result.GetString())
}
