package main

import (
	"log"
	"os"

	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/config"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/internal/database"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/internal/model/dto"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/server"
)

func main() {
	cfg := config.LoadConfig()
	db := database.NewPostgresDatabase(cfg)

	migrateAndSeed(db)

	l := log.New(os.Stdout, "orders-gin ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	server.NewGinServer(cfg, db, l).Start()
}

func migrateAndSeed(db database.Database) {
	if err := db.GetDb().Migrator().AutoMigrate(&dto.OrderDTO{}, &dto.OrderItemDTO{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
