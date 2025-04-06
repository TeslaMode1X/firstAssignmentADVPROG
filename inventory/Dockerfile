FROM golang:1.23.5-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./inventory ./inventory

WORKDIR /app/inventory
RUN go build -o inventory_service ./cmd/inventory/main.go

EXPOSE 8081
CMD ["./inventory_service"]