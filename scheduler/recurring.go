package scheduler

import (
	"log"
	"time"

	"github.com/abilfida/go-flexible-scheduler/database"
	"github.com/abilfida/go-flexible-scheduler/task"
	"gorm.io/gorm"
)

const recurringCheckInterval = 30 * time.Second

// StartRecurringScheduler starts the recurring task scheduler
func StartRecurringScheduler() {
	log.Println("Starting Recurring Task Scheduler...")
	ticker := time.NewTicker(recurringCheckInterval)
	defer ticker.Stop()

	for range ticker.C {
		findAndSpawnRecurringTasks()
	}
}

func findAndSpawnRecurringTasks() {
	log.Println("RecurringScheduler: Checking for due recurring tasks...")

	var recurringTasks []task.RecurringTask
	now := time.Now()

	// Find active recurring tasks that are due to run
	result := database.DB.Where(
		"is_active = ? AND next_run <= ? AND start_date <= ?",
		true, now, now,
	).Find(&recurringTasks)

	if result.Error != nil {
		log.Printf("RecurringScheduler: Error querying recurring tasks: %v", result.Error)
		return
	}

	if len(recurringTasks) == 0 {
		return
	}

	log.Printf("RecurringScheduler: Found %d due recurring tasks", len(recurringTasks))

	// Process each due recurring task
	for _, rt := range recurringTasks {
		processRecurringTask(rt, now)
	}
}

func processRecurringTask(rt task.RecurringTask, now time.Time) {
	// Double-check if task should run (considering end_date, max_runs)
	if !rt.ShouldRun(now) {
		log.Printf("RecurringScheduler: Task %s (ID: %d) should not run, skipping", rt.Name, rt.ID)
		return
	}

	// Create a new pending task from the recurring task template
	pendingTask := task.TaskPending{
		TaskCore: rt.TaskCore,
	}

	// Set scheduled_at to current time for immediate execution
	pendingTask.ScheduledAt = now.Format("2006-01-02 15:04:05")

	// Use transaction to ensure atomicity
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Create the pending task
		if err := tx.Create(&pendingTask).Error; err != nil {
			log.Printf("RecurringScheduler: Failed to create pending task for recurring task %d: %v", rt.ID, err)
			return err
		}

		// Update the recurring task's next run time and counters
		if err := rt.UpdateNextRun(tx); err != nil {
			log.Printf("RecurringScheduler: Failed to update next run for recurring task %d: %v", rt.ID, err)
			return err
		}

		return nil
	})

	if err != nil {
		log.Printf("RecurringScheduler: Transaction failed for recurring task %d: %v", rt.ID, err)
		return
	}

	log.Printf("RecurringScheduler: Successfully spawned task for '%s' (ID: %d), next run: %s", 
		rt.Name, rt.ID, rt.NextRun.Format("2006-01-02 15:04:05"))

	// Check if task should be deactivated (reached max runs or end date)
	if rt.MaxRuns != nil && rt.RunCount >= *rt.MaxRuns {
		rt.IsActive = false
		database.DB.Save(&rt)
		log.Printf("RecurringScheduler: Deactivated recurring task '%s' (ID: %d) - reached max runs (%d)", 
			rt.Name, rt.ID, *rt.MaxRuns)
	} else if rt.EndDate != nil && now.After(*rt.EndDate) {
		rt.IsActive = false
		database.DB.Save(&rt)
		log.Printf("RecurringScheduler: Deactivated recurring task '%s' (ID: %d) - reached end date", 
			rt.Name, rt.ID)
	}
}
