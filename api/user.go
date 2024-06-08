package api

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"test_system/config"
	"test_system/internal"
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

	out, err := os.CreateTemp(config.TestDir, "*."+language)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		return err
	}

	ts := internal.NewRun(out.Name(), language, problem, username)
	result, err := ts.RunTests()
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(TestResultResponse{Message: result.GetString()})
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
	return json.NewEncoder(w).Encode(internal.Contest.Contestants[id.Id].Runs)
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

func InitUserApi() {
	http.Handle("/api/test", logMiddleware(test))
	http.Handle("/api/problems", logMiddleware(getProblems))
	http.Handle("/api/me", logMiddleware(getMe))
	http.Handle("/api/runs", logMiddleware(getRuns))
	http.Handle("/api/languages", logMiddleware(getLanguages))
	http.Handle("/api/results", logMiddleware(getResults))
	http.Handle("/api/startTime", logMiddleware(getContestStartTime))
}
