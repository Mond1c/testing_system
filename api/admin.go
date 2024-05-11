package api

import (
	"strconv"

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

// getRunInfo return the information of the sxpecified run
func getRunInfo(c *fiber.Ctx) error {
	username := c.Query("username")
	id, ok := config.TestConfig.Credentials[username]
	if !ok {
		return c.Status(fiber.StatusBadRequest).SendString("User not found")
	}
	runIDStr := c.Query("run_id")
	runID, err := strconv.Atoi(runIDStr)
	if err != nil || runID < 0 || runID >= len(internal.Contest.Contestants[id.Id].Runs) {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid run ID")
	}
	if runID >= len(internal.Contest.Contestants[id.Id].Runs) {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid run ID")
	}
	run := internal.Contest.Contestants[id.Id].Runs[runID]
	return c.JSON(run)
}

// getAllRuns returns all runs of all users
func getAllRuns(c *fiber.Ctx) error {
	runs := make([]internal.RunInfo, 0)
	for _, contestant := range internal.Contest.Contestants {
		runs = append(runs, contestant.Runs...)
	}
	return c.JSON(runs)
}

// InitAdminAPI initializes the admin API
func InitAdminAPI(app *fiber.App) {
	app.Get("/api/admin/runs", getRunsOfUser)
	app.Get("/api/admin/run", getRunInfo)
	app.Get("/api/admin/all_runs", getAllRuns)
}
