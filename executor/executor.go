package executor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/abilfida/go-flexible-scheduler/database"
	"github.com/abilfida/go-flexible-scheduler/task"
	"gorm.io/gorm"
)

const (
	maxRetries     = 3
	retryInterval  = 5 * time.Second
	requestTimeout = 30 * time.Second
)

// ExecuteTask sekarang menerima TaskRunning
func ExecuteTask(runningTask task.TaskRunning) {
	var lastError error
	var success bool = false
	var resp *http.Response

	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Printf("Executor: Menjalankan task ID %d (Attempt %d/%d)...", runningTask.ID, attempt, maxRetries)
		runningTask.RetryCount = attempt - 1
		database.DB.Model(&runningTask).Update("retry_count", runningTask.RetryCount)

		client := &http.Client{Timeout: requestTimeout}
		req, err := buildRequest(runningTask)
		if err != nil {
			lastError = fmt.Errorf("gagal membuat request: %w", err)
			break
		}

		resp, err = client.Do(req)
		if err != nil {
			lastError = fmt.Errorf("attempt %d gagal: %w", attempt, err)
			if attempt < maxRetries {
				time.Sleep(retryInterval)
			}
			continue
		}

		success = true
		break
	}

	// Memproses hasil dan memindahkan task
	if !success {
		log.Printf("Executor: Task ID %d Gagal setelah %d attempts. Alasan akhir: %v", runningTask.ID, maxRetries, lastError)
		runningTask.ResponseBody = lastError.Error()
		moveTaskToFinalState(runningTask, task.TaskFailed{TaskCore: runningTask.TaskCore})
	} else {
		defer resp.Body.Close()
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			runningTask.ResponseBody = "Gagal membaca response body: " + err.Error()
			moveTaskToFinalState(runningTask, &task.TaskFailed{TaskCore: runningTask.TaskCore})
		} else {
			log.Printf("Executor: Task ID %d Selesai | Status: %d", runningTask.ID, resp.StatusCode)
			runningTask.ResponseStatusCode = resp.StatusCode
			runningTask.ResponseBody = string(respBody)
			moveTaskToFinalState(runningTask, &task.TaskCompleted{TaskCore: runningTask.TaskCore})
		}
	}
}

// Fungsi baru untuk memindahkan task dari 'running' ke state akhir ('completed' atau 'failed')
func moveTaskToFinalState(runningTask task.TaskRunning, finalState interface{}) {
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Create record di tabel tujuan (completed/failed)
		if err := tx.Create(finalState).Error; err != nil {
			return err
		}
		// 2. Delete record dari tabel 'running'
		if err := tx.Delete(&runningTask).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Printf("FATAL: Gagal memindahkan Task ID %d ke state akhir: %v", runningTask.ID, err)
		// Di sini Anda mungkin perlu menambahkan logic alert/pemulihan
	}

	// Kirim Webhook (setelah task berhasil dipindahkan)
	if runningTask.WebhookURL != "" {
		go sendWebhook(runningTask.TaskCore)
	}
}

// func updateTaskAsFailed(t task.Task, statusCode int, reason string) {
// 	log.Printf("Executor: Task ID %d Gagal. Alasan: %s", t.ID, reason)
// 	t.Status = task.StatusFailed
// 	t.ResponseStatusCode = statusCode
// 	t.ResponseBody = reason
// 	database.DB.Save(&t)

// 	if t.WebhookURL != "" {
// 		go sendWebhook(t)
// 	}
// }

// Fungsi helper untuk menambah headers
func addHeaders(req *http.Request, headersJSON string) {
	if headersJSON == "" {
		return
	}
	var headers map[string]string
	if err := json.Unmarshal([]byte(headersJSON), &headers); err == nil {
		for key, val := range headers {
			req.Header.Set(key, val)
		}
	}
	if req.Method == "POST" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
}

// Fungsi helper untuk menambah query params
func addQueryParams(req *http.Request, queryParamsJSON string) {
	if queryParamsJSON == "" {
		return
	}
	var queryParams map[string]string
	if err := json.Unmarshal([]byte(queryParamsJSON), &queryParams); err == nil {
		q := req.URL.Query()
		for key, val := range queryParams {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}
}

func buildRequest(t task.TaskRunning) (*http.Request, error) {
	requestBody := bytes.NewBuffer([]byte(t.Body))
	req, err := http.NewRequest(strings.ToUpper(t.Method), t.URL, requestBody)
	if err != nil {
		return nil, err
	}
	addHeaders(req, t.Headers)
	addQueryParams(req, t.QueryParams)
	return req, nil
}

func sendWebhook(t task.TaskCore) {
	payload, err := json.Marshal(t)
	if err != nil {
		log.Printf("Webhook: Gagal encode payload untuk Task ID %s: %v", t.ProcessID, err)
		return
	}

	resp, err := http.Post(t.WebhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("Webhook: Gagal mengirim ke %s untuk Task ID %s: %v", t.WebhookURL, t.ProcessID, err)
		return
	}
	defer resp.Body.Close()

	log.Printf("Webhook: Berhasil mengirim notifikasi untuk Task ID %s, status response webhook: %s", t.ProcessID, resp.Status)
}
