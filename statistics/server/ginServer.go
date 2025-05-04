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
	app        *gin.Engine
	db         interfaces.Database
	cfg        *config.Config
	log        *log.Logger
	natsClient *nats.Client     // Add this field
	pubSub     *consumer.PubSub // Add this field to store the PubSub instance
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

	if err := s.initializeStatisticsHttpHandler(); err != nil {
		s.log.Fatalf("Failed to initialize statistics handler: %v", err)
	}

	errCh := make(chan error, 1)
	s.pubSub.Start(context.Background(), errCh)

	go func() {
		for err := range errCh {
			s.log.Printf("NATS error: %v", err)
		}
	}()

	serverUrl := fmt.Sprintf(":%s", s.cfg.Server.Port)
	if err := s.app.Run(serverUrl); err != nil {
		s.log.Panic(err)
	}
}

func (s *ginServer) initializeStatisticsHttpHandler() error {
	statisticsRepo := repository.NewStatisticsRepo(s.db)
	statisticsService := service.NewStatisticsService(statisticsRepo)
	statisticsHandler := handler.NewStatisticsHandler(statisticsService)

	// Create NATS client and store it in the server struct
	var err error
	s.natsClient, err = nats.NewClient(context.Background(), []string{"nats_server:4222"}, "", true) // Remove the NKey if not needed
	if err != nil {
		return fmt.Errorf("failed to connect to NATS: %w", err)
	}
	s.log.Println("NATS connection status is", s.natsClient.Conn.Status().String())

	// Create PubSub and store it in the server struct
	s.pubSub = consumer.NewPubSub(s.natsClient)

	productHandler := natsHandler.NewProductHandler(statisticsService)
	s.pubSub.Subscribe(consumer.PubSubSubscriptionConfig{
		Subject: "inventory.product",
		Handler: productHandler.Handler,
	})

	s.pubSub.Subscribe(consumer.PubSubSubscriptionConfig{
		Subject: "orders.order",
		Handler: productHandler.Handler,
	})

	statisticsGroup := s.app.Group("/statistics")
	{
		statisticsGroup.GET("/inventory", statisticsHandler.GetInventoryStatistics)
		statisticsGroup.GET("/orders", statisticsHandler.GetOrdersStatistics)
	}

	return nil
}

// Add a Stop method to properly close resources
func (s *ginServer) Stop() {
	if s.natsClient != nil {
		s.log.Println("Closing NATS connection...")
		s.natsClient.CloseConnect()
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
