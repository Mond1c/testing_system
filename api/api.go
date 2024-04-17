package api

import (
	"io"
	"log"
	"os"
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

	ts := internal.NewRun(header.Filename, language, problem)
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
	for k, _ := range config.TestConfig.TestsInfo {
		problems = append(problems, k)
	}
	err := c.JSON(ResponseProblems{Problems: problems})
	return err
}

func InitApi(app *fiber.App) {
	app.Post("/api/test", test)
	app.Get("/api/problems", getProblems)
}
