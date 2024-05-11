package api

import (
	"encoding/base64"
	"io"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"

	"test_system/config"
	"test_system/internal"
)

// test tests uploading file with source code for correct working
func test(c *fiber.Ctx) error {
	fileName, err := c.FormFile("file")
	if err != nil {
		return err
	}

	language := c.FormValue("language")
	problem := c.FormValue("problem")
	username := c.FormValue("username")

	file, err := fileName.Open()
	if err != nil {
		return nil
	}
	defer file.Close()
	out, err := os.Create(config.TestDir + "/" + fileName.Filename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return err
	}

	ts := internal.NewRun(config.TestDir+"/"+fileName.Filename, language, problem, username)
	result, err := ts.RunTests()
	if err != nil {
		return err
	}
	err = c.JSON(TestResultResponse{Message: result.GetString()})
	return err
}

// getProblems sends problems for the current contest
func getProblems(c *fiber.Ctx) error {
	return c.JSON(internal.Contest.Problems)
}

// getMe sends name of the current user
func getMe(c *fiber.Ctx) error {
	value := strings.Replace(c.GetReqHeaders()["Authorization"][0], "Basic ", "", 1)
	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return err
	}
	username := strings.Split(string(data), ":")[0]
	return c.JSON(ResponseMe{Username: username})
}

// getResults sends results of the current contest
func getResults(c *fiber.Ctx) error {
	return c.JSON(*internal.Contest)
}

// getRuns sends runs for the specified user
func getRuns(c *fiber.Ctx) error {
	id, ok := config.TestConfig.Credentials[c.Query("name")]
	if !ok {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid username")
	}
	return c.JSON(internal.Contest.Contestants[id.Id].Runs)
}

// getLanguages sends json with languages information
func getLanguages(c *fiber.Ctx) error {
	return c.JSON(config.LangaugesConfig.GetLanguages())
}

// getContestStartTime sends json with the contest start time
func getContestStartTime(c *fiber.Ctx) error {
	return c.JSON(StartTimeResponse{
		StartTime: internal.Contest.StartTime.UnixMilli(),
		Duration:  config.TestConfig.Duration,
	})
}

func InitUserApi(app *fiber.App) {
	app.Post("/api/test", test)
	app.Get("/api/problems", getProblems)
	app.Get("/api/me", getMe)
	app.Get("/api/results", getResults)
	app.Get("/api/languages", getLanguages)
	app.Get("/api/runs", getRuns)
	app.Get("/api/monitor", monitor.New())
	app.Get("/api/startTime", getContestStartTime)
}
