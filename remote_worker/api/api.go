package api

import (
	"encoding/json"
	"github.com/Mond1c/testing_system/remote_worker/config"
	"github.com/Mond1c/testing_system/testing/pkg"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func Test(w http.ResponseWriter, r *http.Request) {
	log.Print(1231231)
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	language := r.FormValue("language")
	problem := r.FormValue("problem")
	username := r.FormValue("username")
	testDir := config.MyConfig.TestDir
	duration := r.FormValue("duration")
	startTime := r.FormValue("startTime")
	userId := r.FormValue("userId")
	filename := r.FormValue("filename")

	out, err := os.CreateTemp(testDir, "Solution*."+language)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if language == "java" {
		newFile, err := os.ReadFile(out.Name())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		output := string(newFile)
		index := strings.Index(output, "class")
		index += 6
		if len(output) <= index {
			http.Error(w, "invalid index in java file", http.StatusInternalServerError)
			return
		}

		for !unicode.IsLetter(rune(output[index])) {
			index++
		}
		className := ""
		for index < len(output) && unicode.IsLetter(rune(output[index])) {
			className += string(output[index])
			index++
		}
		output = strings.Replace(output, className, strings.Split(strings.Split(out.Name(), ".")[0], "/")[1], 1)
		_, err = out.WriteAt([]byte(output), 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	convDuration, err := strconv.ParseInt(duration, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ts := pkg.NewTestingSystem(convDuration, startTime,
		config.MyConfig.TestsInfo[problem].Path, config.MyConfig.TestsInfo[problem].TestsCount, username,
		config.MyConfig.CompileFiles[language])
	run := pkg.NewRun(out.Name(), language, problem, username, userId)
	info, err := ts.RunTests(run)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	info.FileName = filename
	err = json.NewEncoder(w).Encode(*info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func InitAPI() {
	http.HandleFunc("/test", Test)
}
