package api

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"
	"test_system/config"
	"test_system/internal"

	"github.com/gofiber/fiber/v2"
)

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
		return err
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

	err = c.JSON(result)
	internal.CheckForErrorAndSendStatusWithLog(c, err, fiber.StatusInternalServerError)
	return nil
}

type ResponseProblems struct {
	Problems []string `json:"problems"`
}

func getProblems(c *fiber.Ctx) error {
	problems := make([]string, 0, len(config.TestConfig.TestsInfo))
	for k := range config.TestConfig.TestsInfo {
		problems = append(problems, k)
	}
	err := c.JSON(ResponseProblems{Problems: problems})
	return err
}

type ResponseMe struct {
	Username string `json:"username"`
}

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

// TODO: what do if file now replacing
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

// Sends json with languages information.
func getLanguages(c *fiber.Ctx) error {
	_ = c.JSON(internal.LangaugesConfig.GetLanguages())
	log.Printf("%v", internal.LangaugesConfig.GetLanguages())
	return nil
}

func InitApi(app *fiber.App) {
	app.Post("/api/test", test)
	app.Get("/api/problems", getProblems)
	app.Get("/api/me", getMe)
	app.Get("/api/results", getResults)
	app.Get("/api/languages", getLanguages)
	app.Get("/api/runs", getRuns)
}
