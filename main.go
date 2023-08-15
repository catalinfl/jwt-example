package main

import (
	"fmt"

	"github.com/catalinfl/gologin/database"
	"github.com/catalinfl/gologin/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	app := fiber.New()

	err := godotenv.Load()

	database.InitDB()

	if err != nil {
		fmt.Println("Error loading")
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"test": "fuck"})
	})

	routes.Login(app)

	routes.Register(app)

	app.Listen(":3000")
}
