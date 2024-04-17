package api

import (
	"io"
	"log"
	"os"
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

func InitApi(app *fiber.App) {
	app.Post("/api/test", test)
}
