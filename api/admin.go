package api

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"test_system/config"
	"test_system/internal"
)

// getRunsOfUser returns all runs of the specified user
func getRunsOfUser(c *fiber.Ctx) error {
	username := c.Query("username")
	id, ok := config.TestConfig.Credentials[username]
	if !ok {
		return c.Status(fiber.StatusBadRequest).SendString("User not found")
	}
	return c.JSON(internal.Contest.Contestants[id.Id].Runs)
}

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
func getRunInfo(c *fiber.Ctx) error {
	runInfo, err := getRunInfoStruct(c.Query("username"), c.Query("run_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(runInfo)
}

func getUsernameById(id string) string {
	for username, info := range config.TestConfig.Credentials {
		if info.Id == id {
			return username
		}
	}
	return ""
}

func getProgrammingLanguageByExtension(path string) string {
	return strings.Split(path, ".")[1]
}

// getAllRuns returns all runs of all users
func getAllRuns(c *fiber.Ctx) error {
	runs := make([]RunInfoResponse, 0)
	for _, contestant := range internal.Contest.Contestants {
		for i, run := range contestant.Runs {
			runs = append(runs, RunInfoResponse{
				Username: getUsernameById(contestant.Id),
				RunID:    i,
				Problem:  run.Problem,
				Result:   run.Result.GetString(),
				Time:     run.Time,
				Language: getProgrammingLanguageByExtension(run.FileName),
			})
		}
	}
	return c.JSON(runs)
}

func getSourceCodeFileOfUser(c *fiber.Ctx) error {
	runInfo, err := getRunInfoStruct(c.Query("username"), c.Query("run_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.Download(runInfo.FileName)
}

func rejudge(c *fiber.Ctx) error {
	runInfo, err := getRunInfoStruct(c.Query("username"), c.Query("run_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	err = internal.RejudgeRun(runInfo)
	return err
}

// InitAdminAPI initializes the admin API
func InitAdminAPI(app *fiber.App) {
	app.Get("/api/admin/runs", getRunsOfUser)
	app.Get("/api/admin/run", getRunInfo)
	app.Get("/api/admin/all_runs", getAllRuns)
	app.Get("/api/admin/source_code", getSourceCodeFileOfUser)
	app.Post("/api/admin/rejudge", rejudge)
}
