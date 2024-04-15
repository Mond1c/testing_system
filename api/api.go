package api

import (
	"fmt"
	"net/http"
)

// TODO: Delete Me
func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello\n")
}

func InitApi() {
	http.HandleFunc("/", hello)
}
