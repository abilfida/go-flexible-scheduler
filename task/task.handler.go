package task

import (
	"errors"

	"github.com/abilfida/go-flexible-scheduler/database"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// getUserIDFromToken mengekstrak user_id dari JWT
func getUserIDFromToken(c *fiber.Ctx) uint {
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))
	return userID
}

// CreateTask: Menambah task baru
// CreateTask sekarang hanya membuat TaskPending
func CreateTask(c *fiber.Ctx) error {
	task := new(TaskPending)
	if err := c.BodyParser(&task.TaskCore); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	task.UserID = getUserIDFromToken(c)

	result := database.DB.Create(&task)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(task)
}

// GetTasks memfilter berdasarkan query status
func GetTasks(c *fiber.Ctx) error {
	userID := getUserIDFromToken(c)
	status := c.Query("status", "pending") // Default ke 'pending' jika tidak ada query

	var result interface{}
	switch status {
	case "running":
		var tasks []TaskRunning
		database.DB.Where("user_id = ?", userID).Find(&tasks)
		result = tasks
	case "completed":
		var tasks []TaskCompleted
		database.DB.Where("user_id = ?", userID).Find(&tasks)
		result = tasks
	case "failed":
		var tasks []TaskFailed
		database.DB.Where("user_id = ?", userID).Find(&tasks)
		result = tasks
	default: // "pending"
		var tasks []TaskPending
		database.DB.Where("user_id = ?", userID).Find(&tasks)
		result = tasks
	}
	return c.JSON(result)
}

// GetTask mencari task berdasarkan status dan ID
func GetTask(c *fiber.Ctx) error {
	userID := getUserIDFromToken(c)
	status := c.Params("status")
	id := c.Params("id")

	var result interface{}
	var err error

	switch status {
	case "pending":
		var task TaskPending
		err = database.DB.Where("user_id = ?", userID).First(&task, id).Error
		result = task
	case "running":
		var task TaskRunning
		err = database.DB.Where("user_id = ?", userID).First(&task, id).Error
		result = task
	case "completed":
		var task TaskCompleted
		err = database.DB.Where("user_id = ?", userID).First(&task, id).Error
		result = task
	case "failed":
		var task TaskFailed
		err = database.DB.Where("user_id = ?", userID).First(&task, id).Error
		result = task
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Status tidak valid"})
	}

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task tidak ditemukan"})
	}
	return c.JSON(result)
}

// UpdateTask: Mengubah task yang ada
func UpdateTask(c *fiber.Ctx) error {
	userID := getUserIDFromToken(c)
	id := c.Params("id")
	var task TaskPending
	if err := database.DB.Where("user_id = ?", userID).First(&task, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task tidak ditemukan"})
	}

	updateData := new(TaskPending)
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
	var task TaskCore
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

// GetTaskByProcessID mencari task di semua tabel berdasarkan ProcessID.
func GetTaskByProcessID(c *fiber.Ctx) error {
	userID := getUserIDFromToken(c)
	processID := c.Params("process_id")

	var err error

	// Cari di tasks_pending
	var pendingTask TaskPending
	err = database.DB.Where("user_id = ? AND process_id = ?", userID, processID).First(&pendingTask).Error
	if err == nil {
		return c.JSON(pendingTask)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Cari di tasks_running
	var runningTask TaskRunning
	err = database.DB.Where("user_id = ? AND process_id = ?", userID, processID).First(&runningTask).Error
	if err == nil {
		return c.JSON(runningTask)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Cari di tasks_completed
	var completedTask TaskCompleted
	err = database.DB.Where("user_id = ? AND process_id = ?", userID, processID).First(&completedTask).Error
	if err == nil {
		return c.JSON(completedTask)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Cari di tasks_failed
	var failedTask TaskFailed
	err = database.DB.Where("user_id = ? AND process_id = ?", userID, processID).First(&failedTask).Error
	if err == nil {
		return c.JSON(failedTask)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Jika tidak ditemukan di semua tabel
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task with process ID not found"})
}
