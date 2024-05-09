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

func main() {
	applicationConfig := config.ParseArgs()

	CheckIfFileExists(applicationConfig.ConfigPath)
	CheckIfFileExists(applicationConfig.LanguagesPath)

	var err error
	config.TestConfig, err = config.ParseConfig(applicationConfig.ConfigPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	config.LangaugesConfig, err = config.ParseLangauges(applicationConfig.LanguagesPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("Names: %v", config.LangaugesConfig.GetLanguages())

	engine := html.New("./frontend/build", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static("/", "./frontend/build")
	app.Use(basicauth.New(basicauth.Config{
		Users: getBasicAuth(),
	}))
	app.Use(pprof.New())
	app.Use(logger.New(logger.Config{
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}\n",
	}))
	app.Use(limiter.New(limiter.Config{
		Max:        30,
		Expiration: 10 * time.Second,
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
	if applicationConfig.Generate {
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
	err = app.Listen(":" + applicationConfig.Port)
	if err != nil {
		log.Fatal(err)
	}
}
