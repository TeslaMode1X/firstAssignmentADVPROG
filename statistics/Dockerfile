FROM golang:1.23.5-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./statistics ./statistics
COPY ./proto ./proto

WORKDIR /app/statistics
RUN go build -o statistics_service ./cmd/statistics/main.go

EXPOSE 50053
CMD ["./statistics_service"]