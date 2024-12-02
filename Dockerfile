FROM golang:1.23.3 AS builder

# Environment variables for static linking and cross-compilation
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app

# Copy dependency files and download modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the binary with optimisations
RUN go build -ldflags="-s -w -X main.version=1.0.0" -o server cmd/main.go && chmod +x server

# Use a minimal base image for production
FROM gcr.io/distroless/static

WORKDIR /

# Copy the binary and required resources
COPY --from=builder /app/server .
COPY --from=builder /app/templates /templates
COPY --from=builder /app/static /static
COPY --from=builder /app/.env .

# Expose the application's port
EXPOSE 8080

# Run the application
CMD ["./server"]
