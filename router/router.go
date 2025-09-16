package router

import (
	"github.com/abilfida/go-flexible-scheduler/auth"
	"github.com/abilfida/go-flexible-scheduler/middleware"
	"github.com/abilfida/go-flexible-scheduler/task"
	"github.com/gofiber/fiber/v2"
)

func SetupTaskRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	api.Post("/signup", auth.SignUp)
	api.Post("/signin", auth.SignIn)

	taskGroup := api.Group("/tasks", middleware.Protected())
	taskGroup.Post("/", task.CreateTask)
	taskGroup.Get("/", task.GetTasks)
	taskGroup.Get("/:id", task.GetTask)
	taskGroup.Put("/:id", task.UpdateTask)
	taskGroup.Delete("/:id", task.DeleteTask)
}
