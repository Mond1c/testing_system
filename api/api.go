package api

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"

	"test_system/config"
	"test_system/internal"
)

// TestResultResponse represents the result of the run
type TestResultResponse struct {
	Message string `json:"message"`
}

// test tests uploading file with source code for correct working
func test(c *fiber.Ctx) error {
	header, err := c.FormFile("file")
	internal.CheckForErrorAndSendStatusWithLog(c, err, fiber.StatusBadRequest)

	language := c.FormValue("language")
	problem := c.FormValue("problem")
	username := c.FormValue("username")

	file, err := header.Open()
	if err != nil {
		log.Print(err)
		return nil
	}
	defer file.Close()

	out, err := os.Create(config.TestDir + "/" + header.Filename)
	internal.CheckForErrorAndSendStatusWithLog(c, err, fiber.StatusInternalServerError)
	defer out.Close()

	_, err = io.Copy(out, file)
	internal.CheckForErrorAndSendStatusWithLog(c, err, fiber.StatusInternalServerError)

	ts := internal.NewRun(config.TestDir+"/"+header.Filename, language, problem, username)
	result, err := ts.RunTests()
	internal.CheckForErrorAndSendStatusWithLog(c, err, fiber.StatusInternalServerError)

	err = c.JSON(TestResultResponse{Message: result.GetString()})
	internal.CheckForErrorAndSendStatusWithLog(c, err, fiber.StatusInternalServerError)
	return nil
}

// ResponseProblems represents json response for GET /api/problems
type ResponseProblems struct {
	Problems []string `json:"problems"`
}

// getProblems sends problems for the current contest
func getProblems(c *fiber.Ctx) error {
	problems := make([]string, 0, len(config.TestConfig.TestsInfo))
	for k := range config.TestConfig.TestsInfo {
		problems = append(problems, k)
	}
	err := c.JSON(ResponseProblems{Problems: problems})
	return err
}

// ResponseMe represents json response for GET /api/me
type ResponseMe struct {
	Username string `json:"username"`
}

// getMe sends name of the current user
func getMe(c *fiber.Ctx) error {
	value := strings.Replace(c.GetReqHeaders()["Authorization"][0], "Basic ", "", 1)
	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		log.Printf("Can't get username from auth: %v", err)
		return err
	}
	username := strings.Split(string(data), ":")[0]
	_ = c.JSON(ResponseMe{
		Username: username,
	})
	return nil
}

// getResults sends results of the current contest
func getResults(c *fiber.Ctx) error {
	data, err := os.ReadFile(config.TestConfig.OutputPath)
	if err != nil {
		log.Printf("Can't read file: %v", err)
		return nil
	}
	var contest internal.ContestInfo
	err = json.Unmarshal(data, &contest)
	if err != nil {
		log.Printf("Can't parse output contest info: %v", err)
		return nil
	}
	_ = c.JSON(contest)
	return nil
}

// getRuns sends runs for the specified user
func getRuns(c *fiber.Ctx) error {
	data, err := os.ReadFile(config.TestConfig.OutputPath)
	if err != nil {
		log.Printf("Can't read file: %v", err)
		return nil
	}
	var contest internal.ContestInfo
	err = json.Unmarshal(data, &contest)
	if err != nil {
		log.Printf("Can't parse output contest info: %v", err)
		return nil
	}
	id, ok := config.TestConfig.Credentials[c.Query("name")]
	if !ok {
		log.Print("can't find contestant with this login")
		return nil
	}
	_ = c.JSON(contest.Contestants[id.Id].Runs)
	return nil
}

// getLanguages sends json with languages information.
func getLanguages(c *fiber.Ctx) error {
	_ = c.JSON(config.LangaugesConfig.GetLanguages())
	log.Printf("%v", config.LangaugesConfig.GetLanguages())
	return nil
}

// getContestInfo sends json with the current contest info
func getContestInfo(c *fiber.Ctx) error {
	_ = c.JSON(internal.Contest)
	return nil
}

// StartTimeResponse represents startTime of the contest and duration of the contest
type StartTimeResponse struct {
	StartTime int64 `json:"startTime"`
	Duration  int64 `json:"duration"`
}

// getContestStartTime sends json with the contest start time
func getContestStartTime(c *fiber.Ctx) error {
	_ = c.JSON(
		StartTimeResponse{
			StartTime: internal.Contest.StartTime.UnixMilli(),
			Duration:  config.TestConfig.Duration,
		},
	)
	return nil
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
	app.Get("/api/contest", getContestInfo)
	app.Get("/api/startTime", getContestStartTime)
}
