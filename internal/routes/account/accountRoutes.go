package accountRoutes

import (
	"github.com/gofiber/fiber/v2"
	accountHandler "github.com/lutungp/julo-test/internal/handler/account"
)

func SetupAccountRoutes(router fiber.Router) {
	router.Post("/init", accountHandler.CreateAccount)
}
