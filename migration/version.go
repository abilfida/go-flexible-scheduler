package migration

import (
	"github.com/abilfida/go-flexible-scheduler/task"
	"github.com/abilfida/go-flexible-scheduler/user"
)

// Dapatkan semua model yang perlu dimigrasi.
func getVersion1Models() []interface{} {
	return []interface{}{
		&user.User{},
		&task.TaskPending{},   // <-- TAMBAHKAN
		&task.TaskRunning{},   // <-- TAMBAHKAN
		&task.TaskCompleted{}, // <-- TAMBAHKAN
		&task.TaskFailed{},    // <-- TAMBAHKAN
	}
}

// Di masa depan, jika ada V2, Anda bisa membuat fungsi baru:
// func getVersion2Models() []interface{} {
//     return []interface{}{&NewModel{}}
// }
