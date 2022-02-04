package walletRoute

import (
	"github.com/gofiber/fiber/v2"
	walletHandler "github.com/lutungp/julo-test/internal/handler/wallet"
	"github.com/lutungp/julo-test/middleware"
)

func SetupWalletRoute(router fiber.Router) {
	note := router.Group("/wallet", middleware.JWTProtected())

	note.Post("/", walletHandler.EnableWalletCustomer)
	note.Get("/", walletHandler.GetBalance)
	note.Post("/deposits", walletHandler.DepositBalance)
	note.Post("/withdrawals", walletHandler.WithDrawls)
	note.Patch("/", walletHandler.DisableWalletCustomer)
}
