package main

import (
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/config"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/database"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/server"
	"log"
	"os"
)

func main() {
	cfg := config.LoadConfig()
	db := database.NewPostgresDatabase(cfg)
	l := log.New(os.Stdout, "inventory-gin", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	server.NewGinServer(cfg, db, l).Start()
}
