package main

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/template/html/v2"

	"test_system/api"
	"test_system/config"
	"test_system/internal"
)

func Render(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{})
}

func CheckIfFileExists(path string) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		log.Fatal(err)
	}
}

func getBasicAuth() map[string]string {
	data := make(map[string]string)
	for k, v := range config.TestConfig.Credentials {
		data[k] = v.Password
	}
	return data
}

func parseConfig[T any](f func(path string) (T, error), path string) T {
	c, err := f(path)
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func initMiddleware(app *fiber.App) {
	app.Use(basicauth.New(basicauth.Config{
		Users: getBasicAuth(),
	}))

	app.Use(pprof.New())

	app.Use(logger.New(logger.Config{
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}\n",
	}))

	app.Use(limiter.New(limiter.Config{
		Max:        30,
		Expiration: time.Second * 10,
	}))
}

func initFrontend(app *fiber.App) {
	app.Static("/", "./frontend/build")

	for _, url := range []string{
		"/",
		"/upload",
		"/results",
		"/runs",
	} {
		app.Get(url, Render)
	}
}

func main() {
	applicationConfig := config.ParseArgs()

	CheckIfFileExists(applicationConfig.ConfigPath)
	CheckIfFileExists(applicationConfig.LanguagesPath)

	config.TestConfig = parseConfig(config.ParseConfig, applicationConfig.ConfigPath)
	config.LangaugesConfig = parseConfig(config.ParseLangauges, applicationConfig.LanguagesPath)

	log.Printf("Names: %v", config.LangaugesConfig.GetLanguages())

	engine := html.New("./frontend/build", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	initMiddleware(app)

    api.InitUserApi(app)

	initFrontend(app)

	if applicationConfig.Generate {
		err := internal.GenerateContestInfo()
		if err != nil {
			return
		}
	}

	go internal.UpdateContestInfo()

	if _, err := os.Stat(config.TestDir); !errors.Is(err, os.ErrNotExist) {
		err = os.RemoveAll(config.TestDir)
		if err != nil {
			log.Fatal(err)
		}
	}

	err := os.Mkdir(config.TestDir, 0750)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(app.Listen(":" + applicationConfig.Port))
}
