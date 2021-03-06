package noteRoutes

import (
	"github.com/gofiber/fiber/v2"
	noteHandler "github.com/lutungp/julo-test/internal/handler/note"
	"github.com/lutungp/julo-test/middleware"
)

func SetupNoteRoutes(router fiber.Router) {
	note := router.Group("/note", middleware.JWTProtected())
	// Create a Note
	note.Post("/", noteHandler.CreateNotes)
	// Read all Notes
	note.Get("/", noteHandler.GetNotes)
	// Read one Note
	note.Get("/:noteId", noteHandler.GetNote)
	// Update one Note
	note.Put("/:noteId", noteHandler.UpdateNote)
	// Delete one Note
	note.Delete("/:noteId", noteHandler.DeleteNote)
}
