package api

import (
	"encoding/base64"
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
		log.Fatal(err)
		return err
	}
	defer file.Close()

	out, err := os.Create(header.Filename)
	internal.CheckForErrorAndSendStatusWithLog(c, err, fiber.StatusInternalServerError)
	defer out.Close()

	_, err = io.Copy(out, file)
	internal.CheckForErrorAndSendStatusWithLog(c, err, fiber.StatusInternalServerError)
	defer internal.RemoveFile(header.Filename)

	ts := internal.NewRun(header.Filename, language, problem, username)
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
		log.Fatalf("Can't get username from auth: %v", err)
		return err
	}
	username := strings.Split(string(data), ":")[0]
	c.JSON(ResponseMe{
		Username: username,
	})
	return nil
}

func InitApi(app *fiber.App) {
	app.Post("/api/test", test)
	app.Get("/api/problems", getProblems)
	app.Get("/api/me", getMe)
}
