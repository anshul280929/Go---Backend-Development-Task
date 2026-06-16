package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/user/user-management-api/internal/handler"
)

// Setup registers all application routes on the Fiber app.
func Setup(app *fiber.App, userHandler *handler.UserHandler) {
	api := app.Group("/users")

	api.Post("/", userHandler.CreateUser)
	api.Get("/", userHandler.ListUsers)
	api.Get("/:id", userHandler.GetUser)
	api.Put("/:id", userHandler.UpdateUser)
	api.Delete("/:id", userHandler.DeleteUser)
}
