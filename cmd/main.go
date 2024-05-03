package main

import (
	"errors"
	"log"
	"os"
	"test_system/api"
	"test_system/config"
	"test_system/internal"

	"flag"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/template/html/v2"
)

func Render(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{})
}

func CheckIfFileExists(path string) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		log.Fatal(err)
	}
}

func main() {
	port := flag.String("flag", "8080", "port for the server")
	configPath := flag.String("config", "", "path to the config file")
	langaugesPath := flag.String("languages", "", "path to the languages settings file")
	generateOutput := flag.Bool("generate", false, "set it if you want generate output json file (turn on on first run)")
	flag.Parse()

	CheckIfFileExists(*configPath)
	CheckIfFileExists(*langaugesPath)

	var err error
	config.TestConfig, err = config.ParseConfig(*configPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	internal.LangaugesConfig, err = internal.ParseLangauges(*langaugesPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("Names: %v", internal.LangaugesConfig.GetLanguages())

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
		"/results",
		"/runs",
	}

	for _, route := range frontendRoutes {
		app.Get(route, Render)
	}
	if *generateOutput {
		err = internal.GenerateContestInfo()
		if err != nil {
			return
		}
	}
	go internal.UpdateContestInfo()
	if _, err = os.Stat(config.TestDir); !errors.Is(err, os.ErrNotExist) {
		err = os.RemoveAll(config.TestDir)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = os.Mkdir(config.TestDir, 0750)
	if err != nil {
		log.Fatal(err)
	}
	err = app.Listen(":" + *port)
	if err != nil {
		log.Fatal(err)
	}
}
