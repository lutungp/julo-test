package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lutungp/julo-test/config"
	"github.com/lutungp/julo-test/database"
	"github.com/lutungp/julo-test/router"
)

func main() {
	// Start a new fiber app
	app := fiber.New()

	// Connect to the Database
	database.ConnectDB()

	// Setup the router
	router.SetupRoutes(app)

	// Listen on PORT 3000
	app.Listen(config.Config("PORT"))
}
