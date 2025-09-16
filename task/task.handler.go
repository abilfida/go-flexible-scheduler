package task

import (
	"github.com/abilfida/go-flexible-scheduler/database"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// getUserIDFromToken mengekstrak user_id dari JWT
func getUserIDFromToken(c *fiber.Ctx) uint {
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))
	return userID
}

// CreateTask: Menambah task baru
func CreateTask(c *fiber.Ctx) error {
	task := new(Task)
	if err := c.BodyParser(task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	userID := getUserIDFromToken(c)
	task.UserID = userID // Set pemilik task
	task.Status = StatusPending

	result := database.DB.Create(&task)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(task)
}

// GetTasks: Menampilkan semua task
func GetTasks(c *fiber.Ctx) error {
	userID := getUserIDFromToken(c)
	var tasks []Task
	database.DB.Where("user_id = ?", userID).Find(&tasks)
	return c.JSON(tasks)
}

// GetTask: Menampilkan satu task berdasarkan ID
func GetTask(c *fiber.Ctx) error {
	userID := getUserIDFromToken(c)
	id := c.Params("id")
	var task Task
	result := database.DB.Where("user_id = ?", userID).First(&task, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task tidak ditemukan"})
	}
	return c.JSON(task)
}

// UpdateTask: Mengubah task yang ada
func UpdateTask(c *fiber.Ctx) error {
	userID := getUserIDFromToken(c)
	id := c.Params("id")
	var task Task
	if err := database.DB.Where("user_id = ?", userID).First(&task, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task tidak ditemukan"})
	}

	updateData := new(Task)
	if err := c.BodyParser(updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	database.DB.Model(&task).Updates(updateData)
	return c.JSON(task)
}

// DeleteTask: Menghapus task
func DeleteTask(c *fiber.Ctx) error {
	userID := getUserIDFromToken(c)
	id := c.Params("id")
	var task Task
	// Verifikasi kepemilikan sebelum menghapus
	if err := database.DB.Where("user_id = ?", userID).First(&task, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task tidak ditemukan"})
	}

	result := database.DB.Delete(&task)
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Gagal menghapus task"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
