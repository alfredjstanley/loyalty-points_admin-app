# Use official Golang image
FROM golang:1.20 AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o server cmd/main.go

# Use minimal image for running the app
FROM gcr.io/distroless/base-debian10

WORKDIR /

# Copy the compiled binary
COPY --from=builder /app/server .

# Copy templates
COPY --from=builder /app/templates ./templates

# Expose the application port
EXPOSE 8080

# Run the server
CMD ["./server"]
