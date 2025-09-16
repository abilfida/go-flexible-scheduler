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
)

const (
	maxRetries     = 3
	retryInterval  = 5 * time.Second
	requestTimeout = 30 * time.Second // Timeout 30 detik
)

func ExecuteTask(t task.Task) {
	log.Printf("Executor: Menjalankan task ID %d -> %s %s", t.ID, t.Method, t.URL)

	var lastError error
	var success bool = false
	var resp *http.Response

	// Loop untuk retry
	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Printf("Executor: Menjalankan task ID %d (Attempt %d/%d)...", t.ID, attempt, maxRetries)

		// Update retry count di DB untuk setiap percobaan
		database.DB.Model(&t).Update("retry_count", attempt-1)

		// 1. Siapkan Request
		client := &http.Client{Timeout: requestTimeout}
		requestBody := bytes.NewBuffer([]byte(t.Body))
		req, err := http.NewRequest(strings.ToUpper(t.Method), t.URL, requestBody)
		if err != nil {
			lastError = fmt.Errorf("gagal membuat request: %w", err)
			log.Printf("Executor: Gagal membuat request untuk Task ID %d: %v", t.ID, err)
			break // Keluar loop jika request tidak valid
		}

		// 2. Tambahkan Headers dan Query Params
		addHeaders(req, t.Headers)
		addQueryParams(req, t.QueryParams)

		// 3. Lakukan HTTP Request
		resp, err = client.Do(req)
		if err != nil {
			lastError = fmt.Errorf("attempt %d gagal: %w", attempt, err)
			log.Printf("Executor: Attempt %d Gagal untuk Task ID %d: %v", attempt, t.ID, err)

			// Tunggu sebelum retry berikutnya jika bukan attempt terakhir
			if attempt < maxRetries {
				time.Sleep(retryInterval)
			}
			continue // Lanjutkan ke attempt berikutnya
		}

		// Jika request berhasil (tidak ada error jaringan/timeout), anggap sukses dan keluar loop
		success = true
		break
	}

	// 4. Proses Hasil Akhir
	if !success {
		// Jika semua attempt gagal
		log.Printf("Executor: Task ID %d Gagal setelah %d attempts. Alasan akhir: %v", t.ID, maxRetries, lastError)
		t.Status = task.StatusFailed
		if lastError != nil {
			t.ResponseBody = lastError.Error()
		}
	} else {
		// Jika salah satu attempt berhasil
		defer resp.Body.Close()
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Status = task.StatusFailed
			t.ResponseBody = "Gagal membaca response body: " + err.Error()
		} else {
			t.Status = task.StatusCompleted
			t.ResponseBody = string(respBody)
			t.ResponseStatusCode = resp.StatusCode
		}
		log.Printf("Executor: Task ID %d Selesai | Status: %d", t.ID, resp.StatusCode)
	}

	// 5. Simpan status akhir ke Database
	database.DB.Save(&t)

	// 6. Kirim Webhook dengan hasil akhir
	if t.WebhookURL != "" {
		go sendWebhook(t)
	}
}

func updateTaskAsFailed(t task.Task, statusCode int, reason string) {
	log.Printf("Executor: Task ID %d Gagal. Alasan: %s", t.ID, reason)
	t.Status = task.StatusFailed
	t.ResponseStatusCode = statusCode
	t.ResponseBody = reason
	database.DB.Save(&t)

	if t.WebhookURL != "" {
		go sendWebhook(t)
	}
}

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

func sendWebhook(t task.Task) {
	payload, err := json.Marshal(t)
	if err != nil {
		log.Printf("Webhook: Gagal encode payload untuk Task ID %d: %v", t.ID, err)
		return
	}

	resp, err := http.Post(t.WebhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("Webhook: Gagal mengirim ke %s untuk Task ID %d: %v", t.WebhookURL, t.ID, err)
		return
	}
	defer resp.Body.Close()

	log.Printf("Webhook: Berhasil mengirim notifikasi untuk Task ID %d, status response webhook: %s", t.ID, resp.Status)
}
