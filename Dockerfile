FROM golang:1.23.3 AS builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o server cmd/main.go && chmod +x server

FROM gcr.io/distroless/static

WORKDIR /

COPY --from=builder /app/server .

COPY --from=builder /app/templates /templates

EXPOSE 8080


CMD ["./server"]
