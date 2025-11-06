package task

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RecurringTaskPattern defines recurring pattern types
type RecurringTaskPattern string

const (
	PatternDaily     RecurringTaskPattern = "daily"      // setiap hari pada waktu tertentu
	PatternHourly    RecurringTaskPattern = "hourly"     // setiap jam pada menit tertentu
	PatternMinutely  RecurringTaskPattern = "minutely"   // setiap X menit
	PatternWeekly    RecurringTaskPattern = "weekly"     // setiap minggu pada hari tertentu
	PatternMonthly   RecurringTaskPattern = "monthly"    // setiap bulan pada tanggal tertentu
	PatternCustom    RecurringTaskPattern = "custom"     // cron-like expression
)

// RecurringTask untuk task yang berjalan berulang
type RecurringTask struct {
	gorm.Model
	
	// Template Task yang akan dijalankan
	TaskCore
	
	// Recurring Configuration
	Pattern     RecurringTaskPattern `json:"pattern" gorm:"not null"`              // Pola pengulangan
	Interval    int                  `json:"interval" gorm:"default:1"`            // Interval untuk pattern minutely/hourly
	Time        string               `json:"time" gorm:"type:varchar(5)"`          // Format HH:MM untuk daily/weekly
	DayOfWeek   int                  `json:"day_of_week" gorm:"default:0"`         // 0=Sunday, 1=Monday, dst untuk weekly
	DayOfMonth  int                  `json:"day_of_month" gorm:"default:1"`        // 1-31 untuk monthly
	CronExpr    string               `json:"cron_expr" gorm:"type:varchar(100)"`   // Custom cron expression
	
	// Control Fields
	IsActive      bool      `json:"is_active" gorm:"default:true"`
	StartDate     time.Time `json:"start_date" gorm:"not null"`
	EndDate       *time.Time `json:"end_date,omitempty"`                      // Optional end date
	LastRun       *time.Time `json:"last_run,omitempty"`
	NextRun       time.Time `json:"next_run" gorm:"not null;index"`
	RunCount      int       `json:"run_count" gorm:"default:0"`
	MaxRuns       *int      `json:"max_runs,omitempty"`                      // Optional limit
	
	// Metadata
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description" gorm:"type:text"`
}

// BeforeCreate hook to generate ProcessID and calculate next run
func (rt *RecurringTask) BeforeCreate(tx *gorm.DB) (err error) {
	if rt.ProcessID == "" {
		rt.ProcessID = uuid.New().String()
	}
	
	// Calculate next run if not set
	if rt.NextRun.IsZero() {
		rt.NextRun = rt.CalculateNextRun(rt.StartDate)
	}
	
	return
}

// TableName for GORM
func (RecurringTask) TableName() string {
	return "recurring_tasks"
}

// CalculateNextRun calculates the next execution time based on pattern
func (rt *RecurringTask) CalculateNextRun(from time.Time) time.Time {
	switch rt.Pattern {
	case PatternDaily:
		return rt.calculateDailyNext(from)
	case PatternHourly:
		return rt.calculateHourlyNext(from)
	case PatternMinutely:
		return rt.calculateMinutelyNext(from)
	case PatternWeekly:
		return rt.calculateWeeklyNext(from)
	case PatternMonthly:
		return rt.calculateMonthlyNext(from)
	case PatternCustom:
		return rt.calculateCustomNext(from)
	default:
		return from.Add(time.Hour) // fallback
	}
}

// calculateDailyNext - setiap hari pada waktu tertentu (misal 07:30)
func (rt *RecurringTask) calculateDailyNext(from time.Time) time.Time {
	hour, minute := rt.parseTime()
	next := time.Date(from.Year(), from.Month(), from.Day(), hour, minute, 0, 0, from.Location())
	
	// Jika waktu hari ini sudah lewat, ambil hari berikutnya
	if next.Before(from) || next.Equal(from) {
		next = next.AddDate(0, 0, 1)
	}
	
	return next
}

// calculateHourlyNext - setiap jam pada menit tertentu (misal menit ke-30)
func (rt *RecurringTask) calculateHourlyNext(from time.Time) time.Time {
	_, minute := rt.parseTime()
	next := time.Date(from.Year(), from.Month(), from.Day(), from.Hour(), minute, 0, 0, from.Location())
	
	// Jika menit ini sudah lewat, ambil jam berikutnya
	if next.Before(from) || next.Equal(from) {
		next = next.Add(time.Hour)
	}
	
	return next
}

// calculateMinutelyNext - setiap X menit
func (rt *RecurringTask) calculateMinutelyNext(from time.Time) time.Time {
	interval := rt.Interval
	if interval <= 0 {
		interval = 5 // default 5 menit
	}
	
	return from.Add(time.Duration(interval) * time.Minute)
}

// calculateWeeklyNext - setiap minggu pada hari tertentu
func (rt *RecurringTask) calculateWeeklyNext(from time.Time) time.Time {
	hour, minute := rt.parseTime()
	currentWeekday := int(from.Weekday())
	targetWeekday := rt.DayOfWeek
	
	// Hitung berapa hari ke depan untuk mencapai target weekday
	daysUntilTarget := (targetWeekday - currentWeekday + 7) % 7
	if daysUntilTarget == 0 {
		// Hari yang sama, cek apakah waktunya sudah lewat
		targetTime := time.Date(from.Year(), from.Month(), from.Day(), hour, minute, 0, 0, from.Location())
		if targetTime.Before(from) || targetTime.Equal(from) {
			daysUntilTarget = 7 // minggu depan
		}
	}
	
	next := from.AddDate(0, 0, daysUntilTarget)
	return time.Date(next.Year(), next.Month(), next.Day(), hour, minute, 0, 0, next.Location())
}

// calculateMonthlyNext - setiap bulan pada tanggal tertentu
func (rt *RecurringTask) calculateMonthlyNext(from time.Time) time.Time {
	hour, minute := rt.parseTime()
	targetDay := rt.DayOfMonth
	if targetDay <= 0 {
		targetDay = 1
	}
	
	// Coba bulan ini dulu
	next := time.Date(from.Year(), from.Month(), targetDay, hour, minute, 0, 0, from.Location())
	
	// Jika tanggal tidak valid (misal 31 Februari) atau sudah lewat, ambil bulan berikutnya
	if next.Day() != targetDay || next.Before(from) || next.Equal(from) {
		next = time.Date(from.Year(), from.Month()+1, targetDay, hour, minute, 0, 0, from.Location())
		// Validasi lagi untuk bulan berikutnya
		if next.Day() != targetDay {
			// Jika masih tidak valid, ambil hari terakhir bulan tersebut
			next = time.Date(next.Year(), next.Month()+1, 0, hour, minute, 0, 0, next.Location())
		}
	}
	
	return next
}

// calculateCustomNext - implementasi sederhana untuk cron expression
// Untuk implementasi lengkap, bisa menggunakan library seperti "github.com/robfig/cron/v3"
func (rt *RecurringTask) calculateCustomNext(from time.Time) time.Time {
	// Implementasi sederhana - bisa diperluas sesuai kebutuhan
	// Untuk sekarang, fallback ke daily
	return rt.calculateDailyNext(from)
}

// parseTime parses time string (HH:MM) into hour and minute
func (rt *RecurringTask) parseTime() (hour, minute int) {
	if rt.Time == "" {
		return 0, 0
	}
	
	// Parse format HH:MM
	if t, err := time.Parse("15:04", rt.Time); err == nil {
		return t.Hour(), t.Minute()
	}
	
	return 0, 0
}

// ShouldRun checks if the recurring task should run now
func (rt *RecurringTask) ShouldRun(now time.Time) bool {
	if !rt.IsActive {
		return false
	}
	
	// Check if before start date
	if now.Before(rt.StartDate) {
		return false
	}
	
	// Check if after end date
	if rt.EndDate != nil && now.After(*rt.EndDate) {
		return false
	}
	
	// Check max runs
	if rt.MaxRuns != nil && rt.RunCount >= *rt.MaxRuns {
		return false
	}
	
	// Check if it's time to run
	return now.After(rt.NextRun) || now.Equal(rt.NextRun)
}

// UpdateNextRun updates the next run time and increments run count
func (rt *RecurringTask) UpdateNextRun(tx *gorm.DB) error {
	now := time.Now()
	rt.LastRun = &now
	rt.RunCount++
	rt.NextRun = rt.CalculateNextRun(now)
	
	return tx.Save(rt).Error
}
