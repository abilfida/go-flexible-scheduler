package task

import "gorm.io/gorm"

type Status string

const (
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
)

type Task struct {
	gorm.Model
	URL                string `json:"url" gorm:"not null"`
	Method             string `json:"method" gorm:"not null;default:'GET'"`
	Headers            string `json:"headers" gorm:"type:text"`      // Store as JSON string
	QueryParams        string `json:"query_params" gorm:"type:text"` // Store as JSON string
	Body               string `json:"body" gorm:"type:text"`         // Store as JSON string or raw text
	ScheduledAt        string `json:"scheduled_at" gorm:"not null"`  // Format: "YYYY-MM-DD HH:MM:SS"
	Status             Status `json:"status" gorm:"default:'pending'"`
	WebhookURL         string `json:"webhook_url"`
	ResponseStatusCode int    `json:"response_status_code"`
	ResponseBody       string `json:"response_body" gorm:"type:text"`
}
