package db

import "os"

type Database struct {
	User         string `env:"DB_USER, required"`
	Password     string `env:"DB_PASSWORD, required"`
	Host         string `env:"DB_HOST, required"`
	Port         string `env:"DB_PORT, required"`
	DriverName   string `env:"DB_DRIVER, required"`
	DatabaseName string `env:"DB_NAME, required"`
	SSLMode      string `env:"DB_SSLMODE, default=disable"`
}

func InitDbConfig() Database {
	return Database{
		User:         os.Getenv("POSTGRES_USER"),
		Password:     os.Getenv("POSTGRES_PASSWORD"),
		Host:         os.Getenv("POSTGRES_HOST"),
		Port:         os.Getenv("POSTGRES_PORT"),
		DatabaseName: os.Getenv("POSTGRES_DB"),
		DriverName:   os.Getenv("POSTGRES_DRIVER"),
		SSLMode:      os.Getenv("POSTGRES_SSL_MODE"),
	}
}
