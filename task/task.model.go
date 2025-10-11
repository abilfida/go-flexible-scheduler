package task

import (
	"github.com/google/uuid" // <-- IMPORT BARU
	"gorm.io/gorm"
)

// TaskCore berisi semua field yang sama untuk setiap tabel task.
type TaskCore struct {
	ProcessID          string `json:"process_id" gorm:"unique;not null;index"`
	UserID             uint   `json:"user_id"`
	URL                string `json:"url" gorm:"not null"`
	Method             string `json:"method" gorm:"not null;default:'GET'"`
	Headers            string `json:"headers" gorm:"type:text"`
	QueryParams        string `json:"query_params" gorm:"type:text"`
	Body               string `json:"body" gorm:"type:text"`
	ScheduledAt        string `json:"scheduled_at" gorm:"not null"`
	WebhookURL         string `json:"webhook_url"`
	ResponseStatusCode int    `json:"response_status_code"`
	ResponseBody       string `json:"response_body" gorm:"type:text"`
	RetryCount         int    `json:"retry_count" gorm:"default:0"`
	RequestDurationMS  int64  `json:"request_duration_ms"`
	WebhookDurationMS  int64  `json:"webhook_duration_ms"`
}

// TaskPending untuk tugas yang belum dijalankan.
type TaskPending struct {
	gorm.Model
	TaskCore
}

// BeforeCreate adalah GORM Hook yang akan dipanggil sebelum task baru dibuat.
// Kita akan menggunakannya untuk generate ProcessID secara otomatis.
func (t *TaskPending) BeforeCreate(tx *gorm.DB) (err error) {
	// Generate UUID baru jika ProcessID masih kosong
	if t.ProcessID == "" {
		t.ProcessID = uuid.New().String()
	}
	return
}

// TableName secara eksplisit memberitahu GORM nama tabelnya.
func (TaskPending) TableName() string {
	return "tasks_pending"
}

// TaskRunning untuk tugas yang sedang dalam proses eksekusi.
type TaskRunning struct {
	gorm.Model
	TaskCore
}

func (TaskRunning) TableName() string {
	return "tasks_running"
}

// TaskCompleted untuk tugas yang berhasil dijalankan.
type TaskCompleted struct {
	gorm.Model
	TaskCore
}

func (TaskCompleted) TableName() string {
	return "tasks_completed"
}

// TaskFailed untuk tugas yang gagal setelah semua percobaan.
type TaskFailed struct {
	gorm.Model
	TaskCore
}

func (TaskFailed) TableName() string {
	return "tasks_failed"
}
