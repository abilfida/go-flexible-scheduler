package router

import (
	"github.com/abilfida/go-flexible-scheduler/auth"
	"github.com/abilfida/go-flexible-scheduler/middleware"
	"github.com/abilfida/go-flexible-scheduler/task"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	// Rute Otentikasi
	api.Post("/signup", auth.SignUp)
	api.Post("/signin", auth.SignIn)

	// Rute Task
	taskGroup := api.Group("/tasks", middleware.Protected())
	taskGroup.Post("/", task.CreateTask)
	taskGroup.Get("/", task.GetTasks) // Contoh: /api/v1/tasks?status=completed

	// Rute baru untuk pencarian berdasarkan ProcessID
	taskGroup.Get("/process/:process_id", task.GetTaskByProcessID)

	// Rute baru sesuai permintaan
	taskGroup.Get("/:status/:id", task.GetTask)

	// Rute lain yang disesuaikan
	taskGroup.Delete("/:status/:id", task.DeleteTask)
	// taskGroup.Put("/:status/:id", task.UpdateTask) // Jika diimplementasikan
}
