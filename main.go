package main

import (
    "log"
    "github.com/gofiber/fiber/v2"
	"github.com/InfamousFreak/Deeptrade/backend/database"
    "github.com/InfamousFreak/Deeptrade/backend/middlewares"
    "github.com/InfamousFreak/Deeptrade/backend/config"
    "github.com/InfamousFreak/Deeptrade/backend/handlers"
    "github.com/InfamousFreak/Deeptrade/backend/routes"

    "github.com/gofiber/fiber/v2/middleware/cors"
    
)

func main() {
    // start a new fiber app
    app := fiber.New()

    jwt := middlewares.NewAuthMiddleware(config.Secret)

    app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Allows all origins
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))


	if err := database.InitDB(); 
    
    err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	app.Get("/", func(c *fiber.Ctx) error {
        err := c.SendString("And the API is UP!")
        return err
    })

    app.Post("/login", handlers.Login)
    app.Get("/protected", jwt, handlers.Protected)
    
    routes.SetupRouter(app)
    // listen on PORT 300
    app.Listen(":3000")
}

