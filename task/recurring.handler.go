package task

import (
	"errors"
	"strconv"
	"time"

	"github.com/abilfida/go-flexible-scheduler/database"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// CreateRecurringTask creates a new recurring task
func CreateRecurringTask(c *fiber.Ctx) error {
	userID := getUserIDFromToken(c)
	
	var req struct {
		TaskCore
		Name        string                `json:"name" validate:"required"`
		Description string                `json:"description"`
		Pattern     RecurringTaskPattern  `json:"pattern" validate:"required"`
		Interval    int                   `json:"interval"`
		Time        string                `json:"time"`
		DayOfWeek   int                   `json:"day_of_week"`
		DayOfMonth  int                   `json:"day_of_month"`
		CronExpr    string                `json:"cron_expr"`
		IsActive    *bool                 `json:"is_active"`
		StartDate   time.Time             `json:"start_date" validate:"required"`
		EndDate     *time.Time            `json:"end_date"`
		MaxRuns     *int                  `json:"max_runs"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	// Validate required fields
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name is required",
		})
	}
	if req.URL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "URL is required",
		})
	}
	if req.Pattern == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Pattern is required",
		})
	}

	// Validate pattern-specific fields
	if err := validatePatternFields(req.Pattern, req.Time, req.Interval, req.DayOfWeek, req.DayOfMonth); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Set defaults
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	// Create recurring task
	recurringTask := RecurringTask{
		TaskCore: TaskCore{
			UserID:      userID,
			URL:         req.URL,
			Method:      req.Method,
			Headers:     req.Headers,
			QueryParams: req.QueryParams,
			Body:        req.Body,
			WebhookURL:  req.WebhookURL,
		},
		Name:        req.Name,
		Description: req.Description,
		Pattern:     req.Pattern,
		Interval:    req.Interval,
		Time:        req.Time,
		DayOfWeek:   req.DayOfWeek,
		DayOfMonth:  req.DayOfMonth,
		CronExpr:    req.CronExpr,
		IsActive:    isActive,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		MaxRuns:     req.MaxRuns,
	}

	// Set default method if empty
	if recurringTask.Method == "" {
		recurringTask.Method = "GET"
	}

	result := database.DB.Create(&recurringTask)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create recurring task: " + result.Error.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(recurringTask)
}

// GetRecurringTasks retrieves all recurring tasks for the authenticated user
func GetRecurringTasks(c *fiber.Ctx) error {
	userID := getUserIDFromToken(c)
	
	var recurringTasks []RecurringTask
	result := database.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&recurringTasks)
	
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve recurring tasks: " + result.Error.Error(),
		})
	}

	return c.JSON(recurringTasks)
}

// GetRecurringTask retrieves a specific recurring task by ID
func GetRecurringTask(c *fiber.Ctx) error {
	userID := getUserIDFromToken(c)
	id := c.Params("id")

	var recurringTask RecurringTask
	result := database.DB.Where("user_id = ? AND id = ?", userID, id).First(&recurringTask)
	
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Recurring task not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve recurring task: " + result.Error.Error(),
		})
	}

	return c.JSON(recurringTask)
}

// UpdateRecurringTask updates an existing recurring task
func UpdateRecurringTask(c *fiber.Ctx) error {
	userID := getUserIDFromToken(c)
	id := c.Params("id")

	// Find existing task
	var recurringTask RecurringTask
	result := database.DB.Where("user_id = ? AND id = ?", userID, id).First(&recurringTask)
	
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Recurring task not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to find recurring task: " + result.Error.Error(),
		})
	}

	// Parse update request
	var updateReq struct {
		TaskCore
		Name        *string               `json:"name"`
		Description *string               `json:"description"`
		Pattern     *RecurringTaskPattern `json:"pattern"`
		Interval    *int                  `json:"interval"`
		Time        *string               `json:"time"`
		DayOfWeek   *int                  `json:"day_of_week"`
		DayOfMonth  *int                  `json:"day_of_month"`
		CronExpr    *string               `json:"cron_expr"`
		IsActive    *bool                 `json:"is_active"`
		StartDate   *time.Time            `json:"start_date"`
		EndDate     *time.Time            `json:"end_date"`
		MaxRuns     *int                  `json:"max_runs"`
	}

	if err := c.BodyParser(&updateReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	// Track if pattern changed to recalculate next run
	patternChanged := false

	// Update fields if provided
	if updateReq.Name != nil {
		recurringTask.Name = *updateReq.Name
	}
	if updateReq.Description != nil {
		recurringTask.Description = *updateReq.Description
	}
	if updateReq.URL != "" {
		recurringTask.URL = updateReq.URL
	}
	if updateReq.Method != "" {
		recurringTask.Method = updateReq.Method
	}
	if updateReq.Headers != "" {
		recurringTask.Headers = updateReq.Headers
	}
	if updateReq.QueryParams != "" {
		recurringTask.QueryParams = updateReq.QueryParams
	}
	if updateReq.Body != "" {
		recurringTask.Body = updateReq.Body
	}
	if updateReq.WebhookURL != "" {
		recurringTask.WebhookURL = updateReq.WebhookURL
	}
	if updateReq.Pattern != nil {
		recurringTask.Pattern = *updateReq.Pattern
		patternChanged = true
	}
	if updateReq.Interval != nil {
		recurringTask.Interval = *updateReq.Interval
		patternChanged = true
	}
	if updateReq.Time != nil {
		recurringTask.Time = *updateReq.Time
		patternChanged = true
	}
	if updateReq.DayOfWeek != nil {
		recurringTask.DayOfWeek = *updateReq.DayOfWeek
		patternChanged = true
	}
	if updateReq.DayOfMonth != nil {
		recurringTask.DayOfMonth = *updateReq.DayOfMonth
		patternChanged = true
	}
	if updateReq.CronExpr != nil {
		recurringTask.CronExpr = *updateReq.CronExpr
		patternChanged = true
	}
	if updateReq.IsActive != nil {
		recurringTask.IsActive = *updateReq.IsActive
	}
	if updateReq.StartDate != nil {
		recurringTask.StartDate = *updateReq.StartDate
		patternChanged = true
	}
	if updateReq.EndDate != nil {
		recurringTask.EndDate = updateReq.EndDate
	}
	if updateReq.MaxRuns != nil {
		recurringTask.MaxRuns = updateReq.MaxRuns
	}

	// Validate pattern fields if pattern changed
	if patternChanged {
		if err := validatePatternFields(recurringTask.Pattern, recurringTask.Time, recurringTask.Interval, recurringTask.DayOfWeek, recurringTask.DayOfMonth); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		
		// Recalculate next run
		recurringTask.NextRun = recurringTask.CalculateNextRun(time.Now())
	}

	// Save updates
	result = database.DB.Save(&recurringTask)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update recurring task: " + result.Error.Error(),
		})
	}

	return c.JSON(recurringTask)
}

// DeleteRecurringTask deletes a recurring task
func DeleteRecurringTask(c *fiber.Ctx) error {
	userID := getUserIDFromToken(c)
	id := c.Params("id")

	// Find and delete the task
	result := database.DB.Where("user_id = ? AND id = ?", userID, id).Delete(&RecurringTask{})
	
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete recurring task: " + result.Error.Error(),
		})
	}
	
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Recurring task not found",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ToggleRecurringTask toggles the active status of a recurring task
func ToggleRecurringTask(c *fiber.Ctx) error {
	userID := getUserIDFromToken(c)
	id := c.Params("id")

	// Find existing task
	var recurringTask RecurringTask
	result := database.DB.Where("user_id = ? AND id = ?", userID, id).First(&recurringTask)
	
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Recurring task not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to find recurring task: " + result.Error.Error(),
		})
	}

	// Toggle active status
	recurringTask.IsActive = !recurringTask.IsActive

	// If activating, recalculate next run
	if recurringTask.IsActive {
		recurringTask.NextRun = recurringTask.CalculateNextRun(time.Now())
	}

	// Save
	result = database.DB.Save(&recurringTask)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to toggle recurring task: " + result.Error.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Recurring task toggled successfully",
		"active": recurringTask.IsActive,
		"task": recurringTask,
	})
}

// validatePatternFields validates pattern-specific required fields
func validatePatternFields(pattern RecurringTaskPattern, timeStr string, interval, dayOfWeek, dayOfMonth int) error {
	switch pattern {
	case PatternDaily:
		if timeStr == "" {
			return errors.New("time is required for daily pattern (format: HH:MM)")
		}
		if !isValidTimeFormat(timeStr) {
			return errors.New("invalid time format, use HH:MM (24-hour format)")
		}
	case PatternHourly:
		if timeStr == "" {
			return errors.New("time is required for hourly pattern (format: 00:MM for minute)")
		}
		if !isValidTimeFormat(timeStr) {
			return errors.New("invalid time format, use 00:MM for hourly pattern")
		}
	case PatternMinutely:
		if interval <= 0 {
			return errors.New("interval must be greater than 0 for minutely pattern")
		}
	case PatternWeekly:
		if timeStr == "" {
			return errors.New("time is required for weekly pattern (format: HH:MM)")
		}
		if !isValidTimeFormat(timeStr) {
			return errors.New("invalid time format, use HH:MM (24-hour format)")
		}
		if dayOfWeek < 0 || dayOfWeek > 6 {
			return errors.New("day_of_week must be between 0 (Sunday) and 6 (Saturday)")
		}
	case PatternMonthly:
		if timeStr == "" {
			return errors.New("time is required for monthly pattern (format: HH:MM)")
		}
		if !isValidTimeFormat(timeStr) {
			return errors.New("invalid time format, use HH:MM (24-hour format)")
		}
		if dayOfMonth < 1 || dayOfMonth > 31 {
			return errors.New("day_of_month must be between 1 and 31")
		}
	case PatternCustom:
		// For custom patterns, we could validate cron expression here
		// For now, just allow it
	default:
		return errors.New("invalid pattern: must be one of daily, hourly, minutely, weekly, monthly, custom")
	}
	return nil
}

// isValidTimeFormat validates HH:MM format
func isValidTimeFormat(timeStr string) bool {
	if len(timeStr) != 5 {
		return false
	}
	if timeStr[2] != ':' {
		return false
	}
	
	hour, err1 := strconv.Atoi(timeStr[0:2])
	minute, err2 := strconv.Atoi(timeStr[3:5])
	
	if err1 != nil || err2 != nil {
		return false
	}
	
	return hour >= 0 && hour <= 23 && minute >= 0 && minute <= 59
}
