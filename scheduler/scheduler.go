package scheduler

import (
	"log"
	"time"

	"github.com/abilfida/go-flexible-scheduler/database"
	"github.com/abilfida/go-flexible-scheduler/executor"
	"github.com/abilfida/go-flexible-scheduler/task"
	"gorm.io/gorm"
)

const checkInterval = 5 * time.Second

func StartScheduler() {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for range ticker.C {
		findAndExecuteTasks()
	}
}

func findAndExecuteTasks() {
	log.Println("Scheduler: Memeriksa tugas di 'tasks_pending'...")

	var tasksToRun []task.TaskPending
	now := time.Now().Format("2006-01-02 15:04:05")

	// 1. Hanya cari dari tabel tasks_pending
	database.DB.Where("scheduled_at <= ?", now).Find(&tasksToRun)

	if len(tasksToRun) == 0 {
		return
	}

	log.Printf("Scheduler: Menemukan %d tugas untuk dijalankan.", len(tasksToRun))

	for _, t := range tasksToRun {
		// 2. Pindahkan task dari 'pending' ke 'running' dalam sebuah transaksi
		runningTask := task.TaskRunning{TaskCore: t.TaskCore}
		runningTask.ID = t.ID // Pertahankan ID yang sama

		err := database.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&runningTask).Error; err != nil {
				return err
			}
			if err := tx.Delete(&t).Error; err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			log.Printf("Scheduler: Gagal memindahkan Task ID %d ke 'running': %v", t.ID, err)
			continue // Lanjut ke task berikutnya jika gagal
		}

		// 3. Jalankan eksekusi di goroutine
		go executor.ExecuteTask(runningTask)
	}
}
