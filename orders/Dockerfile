FROM golang:1.23.5-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./orders ./orders

WORKDIR /app/orders
RUN go build -o orders_service ./cmd/orders/main.go

EXPOSE 8081
CMD ["./orders_service"]