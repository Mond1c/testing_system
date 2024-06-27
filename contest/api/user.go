// Package api contains all api handlers for the application.
package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/Mond1c/testing_system/contest/config"
	"github.com/Mond1c/testing_system/contest/internal"
	"github.com/Mond1c/testing_system/testing/pkg"
)

// test tests uploading file with source code for correct working
func test(w http.ResponseWriter, r *http.Request) error {
	file, _, err := r.FormFile("file")
	if err != nil {
		return err
	}
	defer file.Close()
	language := r.FormValue("language")
	problem := r.FormValue("problem")
	username := r.FormValue("username")

	out, err := os.CreateTemp(config.TestDir, "Solution*."+language)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		return err
	}

	if language == "java" {
		newFile, err := os.ReadFile(out.Name())
		if err != nil {
			return err
		}
		output := string(newFile)
		index := strings.Index(output, "class")
		index += 6
		if len(output) <= index {
			return errors.New("invalid index in java file")
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
			return err
		}
	}
	testPath, count, err := config.TestConfig.GetTestPathForProblem(problem)
	if err != nil {
		return err
	}
	run := pkg.NewRun(out.Name(), language, problem, username, config.TestConfig.Credentials[username].Id)
	task := pkg.CreateRunTask(config.TestConfig.Duration, config.TestConfig.StartTime,
		testPath, count, username, config.LangaugesConfig.CompileFiles[language],
		run, func(info *pkg.RunInfo) {
			internal.AddRun(internal.Contest, info)
		})
	pkg.MyTestingQueue.PushTask(task)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(TestResultResponse{Message: "Waiting..."})
	return err
}

// getProblems sends problems for the current contest
func getProblems(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(internal.Contest.Problems)
}

// getUsername returns username from the header
func getUsername(header string) (string, error) {
	value := strings.Replace(header, "Basic ", "", 1)
	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}
	return strings.Split(string(data), ":")[0], nil
}

// getMe sends name of the current user
func getMe(w http.ResponseWriter, r *http.Request) error {
	username, err := getUsername(r.Header.Get("Authorization"))
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(UsernameResponse{Username: username})
}

// getResults sends results of the current contest
func getResults(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(*internal.Contest)
}

// getRuns sends runs for the specified user
func getRuns(w http.ResponseWriter, r *http.Request) error {
	username, err := getUsername(r.Header.Get("Authorization"))
	if err != nil {
		return err
	}
	id, ok := config.TestConfig.Credentials[username]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	runs := make([]RunInfoResponse, 0)
	for i, run := range internal.Contest.Contestants[id.Id].Runs {
		runs = append(runs, RunInfoResponse{
			Username: getUsernameById(id.Id),
			RunID:    i,
			Problem:  run.Problem,
			Result:   run.Result.String(),
			Time:     run.Time,
			Language: getProgrammingLanguageByExtension(run.FileName),
		})
	}
	return json.NewEncoder(w).Encode(runs)
}

// getLanguages sends json with languages information
func getLanguages(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(config.LangaugesConfig.GetLanguages())
}

// getContestStartTime sends json with the contest start time
func getContestStartTime(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(StartTimeResponse{
		StartTime: internal.Contest.StartTime.UnixMilli(),
		Duration:  config.TestConfig.Duration,
	})
}

// getSourceCode sends source code of the specified run
func getSourceCode(w http.ResponseWriter, r *http.Request) error {
	username, err := getUsername(r.Header.Get("Authorization"))
	if err != nil {
		return err
	}
	id, ok := config.TestConfig.Credentials[username]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}
	runID := r.FormValue("run_id")
	runIDNumber, err := strconv.Atoi(runID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}
	run := internal.Contest.Contestants[id.Id].Runs[runIDNumber]
	http.ServeFile(w, r, run.FileName)
	return nil
}

// InitUserApi initializes user api
func InitUserApi() {
	http.Handle("/api/test", logMiddleware(test))
	http.Handle("/api/problems", logMiddleware(getProblems))
	http.Handle("/api/me", logMiddleware(getMe))
	http.Handle("/api/runs", logMiddleware(getRuns))
	http.Handle("/api/languages", logMiddleware(getLanguages))
	http.Handle("/api/results", logMiddleware(getResults))
	http.Handle("/api/startTime", logMiddleware(getContestStartTime))
	http.Handle("/api/source_code", logMiddleware(getSourceCode))
}
