package migration

import (
	"github.com/abilfida/go-flexible-scheduler/task"
)

// Dapatkan semua model yang perlu dimigrasi untuk versi pertama.
func getVersion1Models() []interface{} {
	return []interface{}{
		&task.Task{},
		// Tambahkan model lain di sini jika ada untuk v1
	}
}

// Di masa depan, jika ada V2, Anda bisa membuat fungsi baru:
// func getVersion2Models() []interface{} {
//     return []interface{}{&NewModel{}}
// }
