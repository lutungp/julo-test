package noteRoutes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	noteHandler "github.com/lutungp/julo-test/internal/handler/note"
	"github.com/lutungp/julo-test/middleware"
)

func SetupNoteRoutes(router fiber.Router) {
	note := router.Group("/note", middleware.JWTProtected())
	// Create a Note
	note.Post("/", noteHandler.CreateNotes)
	// Read all Notes
	note.Get("/", func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)

	})
	// Read one Note
	note.Get("/:noteId", noteHandler.GetNote)
	// Update one Note
	note.Put("/:noteId", noteHandler.UpdateNote)
	// Delete one Note
	note.Delete("/:noteId", noteHandler.DeleteNote)
}
