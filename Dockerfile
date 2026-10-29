# ==========================================
# Tahap 1: Build binary Go
# ==========================================
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install ca-certificates & git jika diperlukan
RUN apk add --no-cache ca-certificates git

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Salin seluruh source code
COPY . .

# Compile binary statis dan ramping
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server .

# ==========================================
# Tahap 2: Runtime image (Alpine minimal)
# ==========================================
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Salin binary yang sudah meng-embed seluruh template dan static file
COPY --from=builder /app/server /app/server

# Salin template & static sebagai fallback tambahan
COPY --from=builder /app/templates /app/templates
COPY --from=builder /app/static /app/static

# Port default yang digunakan server (dapat di-override via env PORT)
ENV PORT=8080

EXPOSE 8080

CMD ["/app/server"]
