# Contoh Penggunaan Recurring Tasks

Berikut adalah berbagai contoh penggunaan fitur recurring tasks dengan presisi sampai menit.

## Contoh 1: Backup Harian

Menjalankan backup setiap hari pukul 02:00 pagi:

```json
{
  "name": "Daily Database Backup",
  "description": "Backup database setiap hari pukul 02:00",
  "url": "https://api.myapp.com/backup/database",
  "method": "POST",
  "headers": "{\"Authorization\":\"Bearer backup-token\",\"Content-Type\":\"application/json\"}",
  "body": "{\"type\":\"full\",\"compress\":true}",
  "pattern": "daily",
  "time": "02:00",
  "start_date": "2025-11-07T00:00:00+07:00",
  "webhook_url": "https://monitoring.myapp.com/backup-status",
  "is_active": true
}
```

## Contoh 2: Health Check Setiap Jam

Menjalankan health check setiap jam pada menit ke-15:

```json
{
  "name": "Hourly Health Check",
  "description": "Cek kesehatan sistem setiap jam",
  "url": "https://api.myapp.com/health",
  "method": "GET",
  "pattern": "hourly",
  "time": "00:15",
  "start_date": "2025-11-07T00:00:00+07:00",
  "webhook_url": "https://monitoring.myapp.com/health-status",
  "is_active": true
}
```

## Contoh 3: Monitoring Setiap 3 Menit

Monitor server resources setiap 3 menit:

```json
{
  "name": "Server Resource Monitor",
  "description": "Monitor CPU, Memory, Disk setiap 3 menit",
  "url": "https://api.myapp.com/monitor/resources",
  "method": "POST",
  "headers": "{\"Content-Type\":\"application/json\"}",
  "body": "{\"metrics\":[\"cpu\",\"memory\",\"disk\"]}",
  "pattern": "minutely",
  "interval": 3,
  "start_date": "2025-11-07T00:00:00+07:00",
  "webhook_url": "https://monitoring.myapp.com/resource-alert",
  "is_active": true
}
```

## Contoh 4: Laporan Mingguan

Generate laporan setiap hari Senin pukul 09:00:

```json
{
  "name": "Weekly Sales Report",
  "description": "Generate laporan penjualan mingguan",
  "url": "https://api.myapp.com/reports/weekly-sales",
  "method": "POST",
  "headers": "{\"Authorization\":\"Bearer report-token\",\"Content-Type\":\"application/json\"}",
  "body": "{\"format\":\"pdf\",\"email_recipients\":[\"manager@company.com\",\"sales@company.com\"]}",
  "pattern": "weekly",
  "time": "09:00",
  "day_of_week": 1,
  "start_date": "2025-11-10T00:00:00+07:00",
  "webhook_url": "https://api.myapp.com/report-status",
  "is_active": true
}
```

## Contoh 5: Invoice Bulanan

Generate invoice setiap tanggal 1 pukul 10:00, maksimal 12 kali (1 tahun):

```json
{
  "name": "Monthly Invoice Generation",
  "description": "Generate invoice bulanan untuk semua customer",
  "url": "https://api.myapp.com/billing/generate-invoices",
  "method": "POST",
  "headers": "{\"Authorization\":\"Bearer billing-token\",\"Content-Type\":\"application/json\"}",
  "body": "{\"invoice_type\":\"monthly\",\"send_email\":true}",
  "pattern": "monthly",
  "time": "10:00",
  "day_of_month": 1,
  "start_date": "2025-12-01T00:00:00+07:00",
  "end_date": "2026-12-01T00:00:00+07:00",
  "max_runs": 12,
  "webhook_url": "https://api.myapp.com/billing-status",
  "is_active": true
}
```

## Contoh 6: Reminder Meeting

Kirim reminder meeting setiap hari Jumat pukul 16:00:

```json
{
  "name": "Weekly Team Meeting Reminder",
  "description": "Kirim reminder untuk team meeting mingguan",
  "url": "https://api.slack.com/api/chat.postMessage",
  "method": "POST",
  "headers": "{\"Authorization\":\"Bearer slack-bot-token\",\"Content-Type\":\"application/json\"}",
  "body": "{\"channel\":\"#general\",\"text\":\"📅 Reminder: Team meeting besok Senin jam 09:00 di meeting room A\"}",
  "pattern": "weekly",
  "time": "16:00",
  "day_of_week": 5,
  "start_date": "2025-11-07T00:00:00+07:00",
  "is_active": true
}
```

## Contoh 7: Cleanup Data Lama

Hapus data temporary setiap hari pukul 03:30:

```json
{
  "name": "Daily Temp Data Cleanup",
  "description": "Hapus file temporary dan cache lama",
  "url": "https://api.myapp.com/maintenance/cleanup",
  "method": "DELETE",
  "headers": "{\"Authorization\":\"Bearer maintenance-token\"}",
  "query_params": "older_than=7d&type=temp,cache",
  "pattern": "daily",
  "time": "03:30",
  "start_date": "2025-11-07T00:00:00+07:00",
  "webhook_url": "https://monitoring.myapp.com/cleanup-status",
  "is_active": true
}
```

## Contoh 8: API Rate Limit Reset

Reset counter API rate limit setiap jam pada menit ke-00:

```json
{
  "name": "API Rate Limit Reset",
  "description": "Reset semua counter rate limit API",
  "url": "https://api.myapp.com/rate-limit/reset",
  "method": "POST",
  "headers": "{\"Authorization\":\"Bearer admin-token\"}",
  "pattern": "hourly",
  "time": "00:00",
  "start_date": "2025-11-07T00:00:00+07:00",
  "is_active": true
}
```

## Contoh 9: Sync Data External API

Sinkronisasi data dengan external API setiap 10 menit:

```json
{
  "name": "External Data Sync",
  "description": "Sync data dari external partner API",
  "url": "https://api.myapp.com/sync/external-data",
  "method": "POST",
  "headers": "{\"Authorization\":\"Bearer sync-token\",\"Content-Type\":\"application/json\"}",
  "body": "{\"sources\":[\"partner_a\",\"partner_b\"],\"incremental\":true}",
  "pattern": "minutely",
  "interval": 10,
  "start_date": "2025-11-07T00:00:00+07:00",
  "webhook_url": "https://monitoring.myapp.com/sync-status",
  "is_active": true
}
```

## Contoh 10: Certificate Renewal Check

Cek SSL certificate expiry setiap hari pukul 08:00:

```json
{
  "name": "SSL Certificate Check",
  "description": "Cek tanggal expire SSL certificate",
  "url": "https://api.myapp.com/security/check-certificates",
  "method": "GET",
  "headers": "{\"Authorization\":\"Bearer security-token\"}",
  "pattern": "daily",
  "time": "08:00",
  "start_date": "2025-11-07T00:00:00+07:00",
  "webhook_url": "https://alerts.myapp.com/certificate-expiry",
  "is_active": true
}
```

## Tips Penggunaan

### 1. Webhook untuk Monitoring
Selalu gunakan `webhook_url` untuk mendapatkan notifikasi status eksekusi task:

```json
{
  "webhook_url": "https://monitoring.myapp.com/task-status"
}
```

### 2. Timezone Awareness
Pastikan `start_date` dan `end_date` menggunakan timezone yang benar:

```json
{
  "start_date": "2025-11-07T00:00:00+07:00",
  "end_date": "2025-12-31T23:59:59+07:00"
}
```

### 3. Batasan Eksekusi
Gunakan `max_runs` untuk membatasi jumlah eksekusi:

```json
{
  "max_runs": 100,
  "end_date": "2025-12-31T23:59:59+07:00"
}
```

### 4. Error Handling
Gunakan webhook untuk menangani error dan retry:

```json
{
  "webhook_url": "https://api.myapp.com/handle-task-error"
}
```

### 5. Testing
Untuk testing, buat task dengan interval pendek dan `max_runs` kecil:

```json
{
  "pattern": "minutely",
  "interval": 1,
  "max_runs": 3,
  "name": "Test Task"
}
```

## Monitoring dan Maintenance

### Melihat Status Recurring Task
```bash
GET /api/v1/recurring
```

### Toggle Active/Inactive
```bash
POST /api/v1/recurring/{id}/toggle
```

### Update Konfigurasi
```bash
PUT /api/v1/recurring/{id}
```

### Hapus Recurring Task
```bash
DELETE /api/v1/recurring/{id}
```
