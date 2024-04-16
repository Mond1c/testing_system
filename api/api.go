package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"test_system/internal"
)

// TODO: Delete Me
func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello\n")
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprintf(w, `<html>
<head>
  <title>GoLang HTTP Fileserver</title>
</head>

<body>

<h2>Upload a file</h2>

<form action="/test" method="post" enctype="multipart/form-data">
  <label for="file">Filename:</label>
  <input type="file" name="file" id="file">
  <br>
  <input type="submit" name="submit" value="Submit">
</form>

</body>
</html>`)
	}
}

func test(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal(err)
		return
	}
	defer file.Close()
	out, err := os.Create(header.Filename)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	log.Print("File upload successful")
	defer os.Remove(header.Filename)
	ts := internal.NewTestSystem(header.Filename)
	result, err := ts.Run()
	if err != nil {
		log.Fatal(err)
		return
	}
	io.WriteString(w, result.GetString())
}

func InitApi() {
	http.HandleFunc("/test", test)
	http.HandleFunc("/", uploadHandler)
}
