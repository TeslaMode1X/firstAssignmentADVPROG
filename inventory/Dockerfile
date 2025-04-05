FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY ../go.mod ../go.sum ./
RUN go mod download
COPY . .
RUN go build -o inventory ./cmd/inventory/main.go
EXPOSE 8081
CMD ["./inventory"]