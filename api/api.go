package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

// TestResultResponse represents the result of the run
type TestResultResponse struct {
	Message string `json:"message"`
}

// ResponseProblems represents json response for GET /api/problems
type ResponseProblems struct {
	Problems []string `json:"problems"`
}

// ResponseMe represents json response for GET /api/me
type ResponseMe struct {
	Username string `json:"username"`
}

// StartTimeResponse represents startTime of the contest and duration of the contest
type StartTimeResponse struct {
	StartTime int64 `json:"startTime"`
	Duration  int64 `json:"duration"`
}

// InitApi inits api for the fiber app
func InitApi(app *fiber.App) {
	app.Post("/api/test", test)
	app.Get("/api/problems", getProblems)
	app.Get("/api/me", getMe)
	app.Get("/api/results", getResults)
	app.Get("/api/languages", getLanguages)
	app.Get("/api/runs", getRuns)
	app.Get("/api/monitor", monitor.New())
	app.Get("/api/startTime", getContestStartTime)
}
