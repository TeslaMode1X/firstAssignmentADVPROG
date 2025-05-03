package main

import (
	"github.com/TeslaMode1X/firstAssignmentADVPROG/statistics/config"
	dbInstance "github.com/TeslaMode1X/firstAssignmentADVPROG/statistics/internal/infrastructure/db"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/statistics/server"
	"log"
	"os"
)

func main() {
	cfg := config.LoadConfig()

	db := dbInstance.NewPostgresDatabase(cfg)

	db.Migrate()

	l := log.New(os.Stdout, "statistics-rpc", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	server.NewGinServer(cfg, db, l).Start()
}
