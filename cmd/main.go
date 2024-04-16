package main

import (
	"fmt"
	"log"
	"net/http"
	"test_system/api"
	"test_system/internal"
)

func generateTests() []*internal.Test {
	tests := make([]*internal.Test, 1000)
	for i := 0; i < 1000; i++ {
		tests[i] = internal.NewTest(fmt.Sprintf("%d %d", i, i+1), fmt.Sprintf("%d", i+i+1))
	}
	return tests
}

func main() {
	//var wg sync.WaitGroup
	//start := time.Now()
	//for i := 0; i < 4; i++ {
	//	wg.Add(1)
	//	go func() {
	//		defer wg.Done()
	//		ctx := internal.NewCodeRunnerContext("cmd/tests/main.cpp", fmt.Sprintf("test%d.out", i))
	//		fmt.Println(i)
	//		result, err := ctx.Test(generateTests())
	//		if err != nil {
	//			log.Fatalf("ERROR #%d: %s\n", i, result.GetString())
	//		}
	//		fmt.Printf("RESULT #%d: %s\n", i, result.GetString())
	//	}()
	//}
	//wg.Wait()
	//log.Printf("Time elapsed: %v\n", time.Since(start))
	api.InitApi()
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
		return
	}
}
