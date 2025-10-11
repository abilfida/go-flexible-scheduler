package scheduler

import (
	"log"
	"time"

	"github.com/abilfida/go-flexible-scheduler/database"
	"github.com/abilfida/go-flexible-scheduler/task"
)

// StartCleanupScheduler memulai goroutine untuk membersihkan soft-deleted tasks.
func StartCleanupScheduler(intervalHours int) {
	log.Printf("Scheduler Pembersihan: Akan berjalan setiap %d jam.", intervalHours)

	// Buat ticker sesuai interval yang dikonfigurasi
	ticker := time.NewTicker(time.Duration(intervalHours) * time.Hour)
	defer ticker.Stop()

	// Lakukan pembersihan pertama kali saat aplikasi dimulai, lalu ikuti ticker
	cleanupSoftDeletedTasks()

	for range ticker.C {
		cleanupSoftDeletedTasks()
	}
}

// cleanupSoftDeletedTasks sekarang membersihkan 'tasks_pending' dan 'tasks_running'.
func cleanupSoftDeletedTasks() {
	log.Println("Scheduler Pembersihan: Memulai proses pembersihan...")

	// 1. Membersihkan tabel tasks_pending
	cleanupTable("tasks_pending", &task.TaskPending{})

	// 2. Membersihkan tabel tasks_running
	cleanupTable("tasks_running", &task.TaskRunning{})
}

// Fungsi helper baru untuk membersihkan tabel secara generik
func cleanupTable(tableName string, model interface{}) {
	result := database.DB.Unscoped().Where("deleted_at IS NOT NULL").Delete(model)

	if result.Error != nil {
		log.Printf("Scheduler Pembersihan: Terjadi error saat membersihkan tabel '%s': %v", tableName, result.Error)
		return
	}

	if result.RowsAffected > 0 {
		log.Printf("Scheduler Pembersihan: Berhasil menghapus %d task secara permanen dari tabel '%s'.", result.RowsAffected, tableName)
	} else {
		log.Printf("Scheduler Pembersihan: Tidak ada task untuk dibersihkan dari tabel '%s' saat ini.", tableName)
	}
}
