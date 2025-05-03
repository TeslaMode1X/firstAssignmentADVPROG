package server

import (
	"context"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/config"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/database"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/repository"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/service"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/usecase"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/pkg/nats"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/pkg/nats/producer"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/proto/gen/inventory"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/proto/gen/promotion"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type grpcServerObject struct {
	server *grpc.Server
	cfg    *config.Config
	db     database.Database
	log    *log.Logger
}

func NewGRPCServer(conf *config.Config, db database.Database, log *log.Logger) Server {
	productRepository := repository.NewProductPostgresRepository(db)
	promoteRepository := repository.NewPromotePostgresRepository(db)

	natsClient, err := nats.NewClient(context.Background(), []string{"nats://nats:4222"}, "", false)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer natsClient.Close()

	inventoryProducer := producer.NewInventoryProducer(natsClient)

	productUseCase := usecase.NewProductUsecaseImpl(productRepository, promoteRepository)
	promotionUseCase := usecase.NewPromoteUsecaseImpl(productRepository, promoteRepository)

	grpcServer := grpc.NewServer()

	reflection.Register(grpcServer)

	inventory.RegisterInventoryServiceServer(grpcServer, service.NewInventoryService(productUseCase, inventoryProducer))
	promotion.RegisterPromotionServiceServer(grpcServer, service.NewPromotionService(promotionUseCase))

	return &grpcServerObject{
		server: grpcServer,
		cfg:    conf,
		db:     db,
		log:    log,
	}
}

func (s *grpcServerObject) Start() {
	port := ":50051"
	if s.cfg.Server.Port != "" {
		port = ":" + s.cfg.Server.Port
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		s.log.Fatalf("Failed to listen on %s: %v", port, err)
	}

	s.log.Printf("Starting inventory gRPC server on %s", port)
	if err := s.server.Serve(lis); err != nil {
		s.log.Fatalf("Failed to serve: %v", err)
	}
}
