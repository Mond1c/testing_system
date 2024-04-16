package api

import (
	"github.com/gofiber/fiber/v2"
	"io"
	"log"
	"os"
	"test_system/internal"
)

// TODO: rewrite
func test(c *fiber.Ctx) error {
	header, err := c.FormFile("file")
	if err != nil {
		c.SendStatus(fiber.StatusBadRequest)
		log.Fatal(err)
		return err
	}
	language := c.FormValue("language")
	file, err := header.Open()
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()
	out, err := os.Create(header.Filename)
	if err != nil {
		c.SendStatus(fiber.StatusInternalServerError)
		log.Fatal(err)
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		c.SendStatus(fiber.StatusInternalServerError)
		log.Fatal(err)
		return err
	}
	log.Print("File upload successful")
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Fatal(err)
		}
	}(header.Filename)
	ts := internal.NewRun(header.Filename, language)
	result, err := ts.RunTests()
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = c.JSON(result)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

type SimpleResponse struct {
	Message string `json:"message"`
}

func simple(c *fiber.Ctx) error {
	log.Println(123)
	err := c.JSON(&SimpleResponse{Message: "ok"})
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func InitApi(app *fiber.App) {
	app.Post("/api/test", test)
	app.Get("/api/simple", simple)
}
