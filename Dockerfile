# --- Tahap 1: Build ---
FROM golang:1.24-alpine AS builder

# Instal paket yang dibutuhkan: ca-certificates
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy go.mod dan go.sum untuk caching dependensi
COPY go.mod go.sum ./

# DITAMBAHKAN: Set GOPROXY ke alternatif global yang andal
ENV GOPROXY=https://goproxy.io,direct

RUN go mod download

# Copy seluruh source code
COPY . .

# Build aplikasi Go menjadi binary yang statis
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o flexible-scheduler .

# --- Tahap 2: Deploy ---
FROM alpine:latest

# Install paket tzdata untuk informasi timezone
RUN apk add --no-cache tzdata ca-certificates

WORKDIR /app

# Copy binary yang sudah di-build dari tahap sebelumnya
COPY --from=builder /app/flexible-scheduler .

# Expose port yang digunakan oleh aplikasi
# EXPOSE 3000

# Perintah untuk menjalankan aplikasi saat container dimulai
CMD ["/app/flexible-scheduler"]