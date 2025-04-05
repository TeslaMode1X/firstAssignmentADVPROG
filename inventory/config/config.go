package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	configDb "github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/config/db"
	configSrv "github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/config/server"
	"github.com/joho/godotenv"
)

var EnvName = "inventory"

type Config struct {
	Server *configSrv.Server
	DB     *configDb.Database
}

func LoadConfig() *Config {
	err := loadDotEnv()
	if err != nil {
		log.Printf("Warning: failed to load .env file: %v. Proceeding with defaults or env vars.", err)
	}

	srv := configSrv.InitServerConfig()
	db := configDb.InitDbConfig()

	return &Config{
		DB:     &db,
		Server: &srv,
	}
}

func loadDotEnv() error {
	filePath := fmt.Sprintf(".env.%s", EnvName)

	if _, err := os.Stat(filePath); err == nil {
		return godotenv.Load(filePath)
	}

	filePath = filepath.Join("..", fmt.Sprintf(".env.%s", EnvName))
	return godotenv.Load(filePath)
}
