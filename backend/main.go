package main

import (
    "github.com/gofiber/fiber/v2"
	"github.com/InfamousFreak/Deeptrade/database"
    "github.com/InfamousFreak/Deeptrade/middlewares"
    "github.com/InfamousFreak/Deeptrade/config"
    "github.com/InfamousFreak/Deeptrade/handlers"
    
)

func main() {
    // start a new fiber app
    app := fiber.New()

    jwt := middlewares.NewAuthMiddleware(config.Secret)

	database.ConnectDB()

	app.Get("/", func(c *fiber.Ctx) error {
        err := c.SendString("And the API is UP!")
        return err
    })

    app.Post("/login", handlers.Login)

    app.Get("/protected", jwt, handlers.Protected)
    // listen on PORT 300
    app.Listen(":3000")
}

