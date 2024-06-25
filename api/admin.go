package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Mond1c/test_system/config"
	"github.com/Mond1c/test_system/internal"
)

// getRunsOfUser returns all runs of the specified user
func getRunsOfUser(w http.ResponseWriter, r *http.Request) error {
	username := r.URL.Query().Get("username")
	id, ok := config.TestConfig.Credentials[username]
	if !ok {
		return errors.New("user not found")
	}
	return json.NewEncoder(w).Encode(internal.Contest.Contestants[id.Id].Runs)
}

// getRunInfoStruct returns internal.RunInfo for the specified username and id of the run
func getRunInfoStruct(username, runIDStr string) (*internal.RunInfo, error) {
	id, ok := config.TestConfig.Credentials[username]
	if !ok {
		return nil, errors.New("user not found")
	}
	runID, err := strconv.Atoi(runIDStr)
	if err != nil || runID < 0 || runID >= len(internal.Contest.Contestants[id.Id].Runs) {
		return nil, errors.New("invalid run ID")
	}
	return &internal.Contest.Contestants[id.Id].Runs[runID], nil
}

// getRunInfo return the information of the sxpecified run
func getRunInfo(w http.ResponseWriter, r *http.Request) error {
	runInfo, err := getRunInfoStruct(r.URL.Query().Get("username"), r.URL.Query().Get("run_id"))
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(runInfo)
}

// getUsername returns username for the specified contestant id
func getUsernameById(id string) string {
	for username, info := range config.TestConfig.Credentials {
		if info.Id == id {
			return username
		}
	}
	return ""
}

// getProgrammingLanguageByExtension returns programming language for the specified file
func getProgrammingLanguageByExtension(path string) string {
	return strings.Split(path, ".")[1]
}

// getAllRuns returns all runs of all users
func getAllRuns(w http.ResponseWriter, r *http.Request) error {
	runs := make([]RunInfoResponse, 0)
	for _, contestant := range internal.Contest.Contestants {
		for i, run := range contestant.Runs {
			runs = append(runs, RunInfoResponse{
				Username: getUsernameById(contestant.Id),
				RunID:    i,
				Problem:  run.Problem,
				Result:   run.Result.String(),
				Time:     run.Time,
				Language: getProgrammingLanguageByExtension(run.FileName),
			})
		}
	}
	return json.NewEncoder(w).Encode(runs)
}

// getSourceCodeFileOfUser gets source code of the problem for the specified user
func getSourceCodeFileOfUser(w http.ResponseWriter, r *http.Request) error {
	runInfo, err := getRunInfoStruct(r.URL.Query().Get("username"), r.URL.Query().Get("run_id"))
	if err != nil {
		return err
	}
	http.ServeFile(w, r, runInfo.FileName)
	return nil
}

// rejudge reruns tests for the specified run of the specified user
func rejudge(w http.ResponseWriter, r *http.Request) error {
	runInfo, err := getRunInfoStruct(r.URL.Query().Get("username"), r.URL.Query().Get("run_id"))
	if err != nil {
		return err
	}
	return internal.RejudgeRun(runInfo)
}

// InitAdminAPI initializes the admin API
func InitAdminAPI() {
	http.Handle("/api/admin/runs", logMiddleware(getRunsOfUser))
	http.Handle("/api/admin/run", logMiddleware(getRunInfo))
	http.Handle("/api/admin/all_runs", logMiddleware(getAllRuns))
	http.Handle("/api/admin/source_code", logMiddleware(getSourceCodeFileOfUser))
	http.Handle("/api/admin/rejudge", logMiddleware(rejudge))
}
