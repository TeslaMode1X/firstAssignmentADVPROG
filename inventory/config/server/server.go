package server

import (
	"os"
	"time"
)

type Server struct {
	Addr         string        `env:"SERVER_ADDR, default=localhost"`
	Host         string        `env:"SERVER_HOST, default=localhost"`
	Port         string        `env:"SERVER_PORT, default=8080"`
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT, default=5s"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT, default=10s"`
	IdleTimeout  time.Duration `env:"SERVER_IDLE_TIMEOUT, default=60s"`
}

func InitServerConfig() Server {
	return Server{
		Addr: os.Getenv("SERVER_ADDR"),
		Port: os.Getenv("SERVER_PORT"),
		Host: os.Getenv("SERVER_HOST"),
	}
}
