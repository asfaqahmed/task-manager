package routes

import (
	"task-manager/handlers"
	"task-manager/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/register", handlers.Register)
	api.Post("/login", handlers.Login)

	protected := api.Group("/")
	protected.Use(middleware.JWTMiddleware)

	protected.Get("/tasks", handlers.GetTasks)
	protected.Post("/tasks", handlers.CreateTask)
	protected.Put("/tasks/:id", handlers.UpdateTask)
	protected.Delete("/tasks/:id", handlers.DeleteTask)
}
