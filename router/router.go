package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	accountRoutes "github.com/lutungp/julo-test/internal/routes/account"
	noteRoutes "github.com/lutungp/julo-test/internal/routes/note"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1", logger.New())

	noteRoutes.SetupNoteRoutes(api)
	accountRoutes.SetupAccountRoutes(api)
}
