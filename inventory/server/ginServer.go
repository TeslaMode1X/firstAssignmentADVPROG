package server

import (
	"fmt"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/handler"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/repository"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/usecase"
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
	s.app.Use(gin.Recovery())
	s.app.Use(gin.Logger())

	s.app.GET("/v1/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	s.initializeProductHttpHandler()

	serverUrl := fmt.Sprintf(":%s", s.cfg.Server.Port)
	if err := s.app.Run(serverUrl); err != nil {
		s.log.Panic(err)
	}
}

func (s *ginServer) initializeProductHttpHandler() {
	productRepository := repository.NewProductPostgresRepository(s.db)
	productUseCase := usecase.NewProductUsecaseImpl(productRepository)
	productHandler := handler.NewProductHttpHandler(productUseCase)

	productRouters := s.app.Group("/product")
	{
		productRouters.POST("/create", productHandler.CreateProduct)
		productRouters.GET("", productHandler.GetProducts)
		productRouters.GET(":id", productHandler.GetProductByID)
		productRouters.PATCH("", productHandler.UpdateProduct)
		productRouters.DELETE(":id", productHandler.DeleteProduct)
	}
}
