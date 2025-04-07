package server

import (
	"fmt"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/config"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/internal/database"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/internal/handler"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/internal/handler/client"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/internal/repository"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/internal/usecase"
	"github.com/gin-gonic/gin"
	"log"

	"net/http"
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
	s.app.Use(gin.Recovery())
	s.app.Use(gin.Logger())

	s.app.GET("/v1/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	s.initializeOrderHttpHandler()

	serverUrl := fmt.Sprintf(":%s", s.cfg.Server.Port)
	if err := s.app.Run(serverUrl); err != nil {
		s.log.Panic(err)
	}
}

func (s *ginServer) initializeOrderHttpHandler() {
	orderRepository := repository.NewOrderPostgresRepository(s.db)
	clientRepo := client.NewInventoryClient("http://api_gateway:8080")
	orderUseCase := usecase.NewOrderUsecaseImpl(orderRepository, clientRepo)
	orderHandler := handler.NewOrderHttpHandler(orderUseCase)

	orderRouters := s.app.Group("/order")
	{
		orderRouters.POST("/create", orderHandler.CreateOrder)
		orderRouters.GET("", orderHandler.GetOrders)
		orderRouters.GET(":id", orderHandler.GetOrderByID)
		orderRouters.PATCH("", orderHandler.UpdateOrderStatus)
	}
}
