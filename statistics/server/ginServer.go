package server

import (
	"context"
	"fmt"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/statistics/config"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/statistics/internal/handler"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/statistics/internal/interfaces"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/statistics/internal/repository"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/statistics/internal/service"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/statistics/pkg/nats"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/statistics/pkg/nats/consumer"
	natsHandler "github.com/TeslaMode1X/firstAssignmentADVPROG/statistics/pkg/nats/handler"
	"github.com/gin-gonic/gin"
	"log"

	"net/http"
)

type ginServer struct {
	app *gin.Engine
	db  interfaces.Database
	cfg *config.Config
	log *log.Logger
}

func NewGinServer(conf *config.Config, db interfaces.Database, log *log.Logger) Server {
	ginApp := gin.Default()

	return &ginServer{
		app: ginApp,
		db:  db,
		cfg: conf,
		log: log,
	}
}

func (s *ginServer) Start() {
	s.app.Use(gin.Recovery())
	s.app.Use(gin.Logger())

	s.app.GET("/v1/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	s.initializeStatisticsHttpHandler()

	serverUrl := fmt.Sprintf(":%s", s.cfg.Server.Port)
	if err := s.app.Run(serverUrl); err != nil {
		s.log.Panic(err)
	}
}

func (s *ginServer) initializeStatisticsHttpHandler() {
	statisticsRepo := repository.NewStatisticsRepo(s.db)
	statisticsService := service.NewStatisticsService(statisticsRepo)
	statisticsHandler := handler.NewStatisticsHandler(statisticsService)

	natsClient, err := nats.NewClient(context.Background(), []string{"localhost:4222"}, "SUACSSL3UAHUDXKFSNVUZRF5UHPMWZ6BFDTJ7M6USDXIEDNPPQYYYCU3VY", true)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("NATS connection status is", natsClient.Conn.Status().String())

	pubSub := consumer.NewPubSub(natsClient)

	errCh := make(chan error, 1)
	pubSub.Start(context.Background(), errCh)

	productHandler := natsHandler.NewProductHandler(statisticsService)
	pubSub.Subscribe(consumer.PubSubSubscriptionConfig{
		Subject: "ap2.service_scv.event.created",
		Handler: productHandler.Handler,
	})

	statisticsGroup := s.app.Group("/statistics")
	{
		statisticsGroup.GET("/inventory", statisticsHandler.GetInventoryStatistics)
		statisticsGroup.GET("/orders", statisticsHandler.GetOrdersStatistics)
	}
}

//Nats struct {
//Hosts        []string `env:"NATS_HOSTS,notEmpty" envSeparator:","`
//NKey         string   `env:"NATS_NKEY,notEmpty"`
//IsTest       bool     `env:"NATS_IS_TEST,notEmpty" envDefault:"true"`
//NatsSubjects NatsSubjects
//}

//# NATS
//NATS_HOSTS="localhost:4222,localhost:4222,localhost:4222"
//NATS_NKEY=SUACSSL3UAHUDXKFSNVUZRF5UHPMWZ6BFDTJ7M6USDXIEDNPPQYYYCU3VY #def key

//# NATS Subjects
//NATS_CLIENT_EVENT_SUBJECT=ap2.service_scv.event.created
