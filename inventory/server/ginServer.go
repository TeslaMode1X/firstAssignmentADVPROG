package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/config"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/database"
	"github.com/gin-gonic/gin"
)

type ginServer struct {
	app *gin.Engine
	db  database.Database
	cfg *config.Config
	log *log.Logger
}

func NewGinServer(conf *config.Config, db database.Database, log *log.Logger) Server {
	ginApp := gin.Default()

	return &ginServer{
		app: ginApp,
		db:  db,
		cfg: conf,
		log: log,
	}
}

func (s *ginServer) Start() {
	s.app.GET("/v1/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	serverUrl := fmt.Sprintf(":%s", s.cfg.Server.Port)

	if err := s.app.Run(serverUrl); err != nil {
		s.log.Panic(err)
	}
}
