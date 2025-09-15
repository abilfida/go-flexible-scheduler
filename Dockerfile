# --- Tahap 1: Build ---
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod dan go.sum untuk caching dependensi
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Build aplikasi Go menjadi binary yang statis
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /scheduler-app .

# --- Tahap 2: Deploy ---
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary yang sudah di-build dari tahap sebelumnya
COPY --from=builder /scheduler-app .

# Expose port yang digunakan oleh aplikasi
EXPOSE 3000

# Perintah untuk menjalankan aplikasi saat container dimulai
CMD ["./scheduler-app"]