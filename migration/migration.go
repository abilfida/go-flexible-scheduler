package migration

import (
	"log"

	"gorm.io/gorm"
)

// AutoMigrate akan menjalankan migrasi GORM untuk semua model yang terdaftar.
func AutoMigrate(db *gorm.DB) {
	log.Println("Memulai proses migrasi database...")

	// Jalankan migrasi untuk Versi 1
	errV1 := db.AutoMigrate(getVersion1Models()...)
	if errV1 != nil {
		log.Fatalf("Gagal migrasi database V1: %v", errV1)
	}

	// CONTOH UNTUK MASA DEPAN (Update Versi Selanjutnya)
	// Jika nanti ada V2, Anda tinggal uncomment dan tambahkan di sini.
	// if !db.Migrator().HasTable(&NewModel{}) {
	//     errV2 := db.AutoMigrator().AutoMigrate(getVersion2Models()...)
	//     if errV2 != nil {
	//         log.Fatalf("Gagal migrasi database V2: %v", errV2)
	//     }
	//     log.Println("Migrasi V2 berhasil.")
	// }

	log.Println("Migrasi database berhasil diselesaikan.")
}
