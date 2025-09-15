package executor

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/abilfida/go-flexible-scheduler/database"
	"github.com/abilfida/go-flexible-scheduler/task"
)

func ExecuteTask(t task.Task) {
	log.Printf("Executor: Menjalankan task ID %d -> %s %s", t.ID, t.Method, t.URL)

	var req *http.Request
	var err error

	// 1. Siapkan Request
	client := &http.Client{Timeout: 30 * time.Second}
	requestBody := bytes.NewBuffer([]byte(t.Body))

	req, err = http.NewRequest(strings.ToUpper(t.Method), t.URL, requestBody)
	if err != nil {
		updateTaskAsFailed(t, 0, "Gagal membuat request: "+err.Error())
		return
	}

	// 2. Tambahkan Headers
	if t.Headers != "" {
		var headers map[string]string
		if err := json.Unmarshal([]byte(t.Headers), &headers); err == nil {
			for key, val := range headers {
				req.Header.Set(key, val)
			}
		}
	}
	// Default content-type untuk POST
	if strings.ToUpper(t.Method) == "POST" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// 3. Tambahkan Query Params
	if t.QueryParams != "" {
		var queryParams map[string]string
		if err := json.Unmarshal([]byte(t.QueryParams), &queryParams); err == nil {
			q := req.URL.Query()
			for key, val := range queryParams {
				q.Add(key, val)
			}
			req.URL.RawQuery = q.Encode()
		}
	}

	// 4. Lakukan HTTP Request
	resp, err := client.Do(req)
	if err != nil {
		updateTaskAsFailed(t, 0, "Gagal melakukan request: "+err.Error())
		return
	}
	defer resp.Body.Close()

	// 5. Baca Response Body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		updateTaskAsFailed(t, resp.StatusCode, "Gagal membaca response body: "+err.Error())
		return
	}

	// 6. Update Task sebagai Berhasil
	t.Status = task.StatusCompleted
	t.ResponseStatusCode = resp.StatusCode
	t.ResponseBody = string(respBody)
	database.DB.Save(&t)

	log.Printf("Executor: Task ID %d Selesai | Status: %d", t.ID, resp.StatusCode)

	// 7. Kirim Webhook (jika ada)
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
