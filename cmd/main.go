package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/template/html/v2"
	"test_system/api"
)

func Home(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{})
}

func main() {
	engine := html.New("./frontend/build", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static("/", "./frontend/build")

	api.InitApi(app)

	frontendRoutes := []string{
		"/",
		"/upload",
	}

	for _, route := range frontendRoutes {
		app.Get(route, Home)
	}

	err := app.Listen(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
