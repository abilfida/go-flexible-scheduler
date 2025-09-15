package task

import "github.com/gofiber/fiber/v2"

func SetupTaskRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	api.Post("/tasks", CreateTask)
	api.Get("/tasks", GetTasks)
	api.Get("/tasks/:id", GetTask)
	api.Put("/tasks/:id", UpdateTask)
	api.Delete("/tasks/:id", DeleteTask)
}
