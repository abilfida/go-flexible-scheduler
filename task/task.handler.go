package task

import (
	"github.com/abilfida/go-flexible-scheduler/database"
	"github.com/gofiber/fiber/v2"
)

// CreateTask: Menambah task baru
func CreateTask(c *fiber.Ctx) error {
	task := new(Task)
	if err := c.BodyParser(task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	task.Status = StatusPending
	result := database.DB.Create(&task)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(task)
}

// GetTasks: Menampilkan semua task
func GetTasks(c *fiber.Ctx) error {
	var tasks []Task
	database.DB.Find(&tasks)
	return c.JSON(tasks)
}

// GetTask: Menampilkan satu task berdasarkan ID
func GetTask(c *fiber.Ctx) error {
	id := c.Params("id")
	var task Task
	result := database.DB.First(&task, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task tidak ditemukan"})
	}
	return c.JSON(task)
}

// UpdateTask: Mengubah task yang ada
func UpdateTask(c *fiber.Ctx) error {
	id := c.Params("id")
	var task Task
	if err := database.DB.First(&task, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task tidak ditemukan"})
	}

	updateData := new(Task)
	if err := c.BodyParser(updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Hanya update field yang di-provide
	database.DB.Model(&task).Updates(updateData)
	return c.JSON(task)
}

// DeleteTask: Menghapus task
func DeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	result := database.DB.Delete(&Task{}, id)
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task tidak ditemukan"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
