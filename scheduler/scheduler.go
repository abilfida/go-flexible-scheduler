package scheduler

import (
	"log"
	"time"

	"github.com/abilfida/go-flexible-scheduler/database"
	"github.com/abilfida/go-flexible-scheduler/executor"
	"github.com/abilfida/go-flexible-scheduler/task"
)

const checkInterval = 5 * time.Second // Seberapa sering scheduler memeriksa DB

func StartScheduler() {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for range ticker.C {
		findAndExecuteTasks()
	}
}

func findAndExecuteTasks() {
	log.Println("Scheduler: Memeriksa tugas yang jatuh tempo...")

	var tasksToRun []task.Task
	now := time.Now().Format("2006-01-02 15:04:05")

	// Cari task yang statusnya 'pending' dan waktunya sudah lewat
	database.DB.Where("status = ? AND scheduled_at <= ?", task.StatusPending, now).Find(&tasksToRun)

	if len(tasksToRun) == 0 {
		return
	}

	log.Printf("Scheduler: Menemukan %d tugas untuk dijalankan.", len(tasksToRun))

	for _, t := range tasksToRun {
		// Update status menjadi 'running' agar tidak diambil oleh worker lain
		database.DB.Model(&t).Update("status", task.StatusRunning)

		// Jalankan eksekusi di goroutine terpisah agar tidak memblokir scheduler
		go executor.ExecuteTask(t)
	}
}
