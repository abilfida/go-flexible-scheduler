package main

import (
	"log"

	"github.com/abilfida/go-flexible-scheduler/config"
	"github.com/abilfida/go-flexible-scheduler/database"
	"github.com/abilfida/go-flexible-scheduler/migration"
	"github.com/abilfida/go-flexible-scheduler/scheduler"
	"github.com/abilfida/go-flexible-scheduler/task"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// 1. Muat Konfigurasi
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Tidak dapat memuat konfigurasi: %v", err)
	}

	// 2. Hubungkan ke Database
	database.ConnectDB(cfg.DSN)
	log.Println("Koneksi Database Berhasil.")

	// 3. Jalankan Migrasi Database
	migration.AutoMigrate(database.DB)

	// 4. Lakukan Auto-Migration untuk tabel Task
	database.DB.AutoMigrate(&task.Task{})
	log.Println("Database Migrated.")

	// 5. Inisialisasi Fiber App
	app := fiber.New()

	// 6. Setup Routing
	task.SetupTaskRoutes(app)

	// 7. Jalankan Scheduler di background
	go scheduler.StartScheduler()
	log.Println("Scheduler Engine Dimulai.")

	// 8. Jalankan Server API
	log.Fatal(app.Listen(":" + cfg.Port))
}
