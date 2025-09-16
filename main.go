package main

import (
	"log"
	"time"

	"github.com/abilfida/go-flexible-scheduler/config"
	"github.com/abilfida/go-flexible-scheduler/database"
	"github.com/abilfida/go-flexible-scheduler/migration"
	"github.com/abilfida/go-flexible-scheduler/scheduler"
	"github.com/abilfida/go-flexible-scheduler/task"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Muat Konfigurasi
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Tidak dapat memuat konfigurasi: %v", err)
	}

	// Atur Timezone Aplikasi
	loc, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		log.Fatalf("Gagal memuat timezone: %v", err)
	}
	time.Local = loc // Mengatur timezone default untuk seluruh aplikasi
	log.Printf("Timezone aplikasi diatur ke: %s", loc.String())

	// Hubungkan ke Database
	database.ConnectDB(cfg.DSN)
	log.Println("Koneksi Database Berhasil.")

	// Jalankan Migrasi Database
	migration.AutoMigrate(database.DB)

	// Lakukan Auto-Migration untuk tabel Task
	database.DB.AutoMigrate(&task.Task{})
	log.Println("Database Migrated.")

	// Inisialisasi Fiber App
	app := fiber.New()

	// Setup Routing
	task.SetupTaskRoutes(app)

	// Jalankan Scheduler di background
	go scheduler.StartScheduler()
	log.Println("Scheduler Engine Dimulai.")

	// Jalankan Server API
	log.Fatal(app.Listen(":" + cfg.Port))
}
