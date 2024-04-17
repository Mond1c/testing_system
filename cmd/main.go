package main

import (
	"errors"
	"os"
	"test_system/api"
	"test_system/config"
	"test_system/internal"

	"flag"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/template/html/v2"
)

func Render(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{})
}

func main() {
	port := flag.String("flag", "8080", "port for the server")
	configPath := flag.String("config", "", "path to the config file")
	flag.Parse()

	if _, err := os.Stat(*configPath); errors.Is(err, os.ErrNotExist) {
		log.Fatal(err)
		return
	}
	var err error
	config.TestConfig, err = config.ParseConfig(*configPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	engine := html.New("./frontend/build", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static("/", "./frontend/build")
	app.Use(basicauth.New(basicauth.Config{
		Users: internal.LoginPassword,
	}))
	api.InitApi(app)

	frontendRoutes := []string{
		"/",
		"/upload",
	}

	for _, route := range frontendRoutes {
		app.Get(route, Render)
	}
	err = internal.GenerateContestInfo()
	if err != nil {
		return
	}
	go internal.UpdateContestInfo()

	err = app.Listen(":" + *port)
	if err != nil {
		log.Fatal(err)
	}
}
