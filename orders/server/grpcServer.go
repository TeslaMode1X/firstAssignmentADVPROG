package server

import (
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/config"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/internal/database"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/internal/handler/client"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/internal/repository"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/internal/service"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/internal/usecase"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/proto/gen/orders"
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
	orderRepository := repository.NewOrderPostgresRepository(db)
	clientRepo := client.NewInventoryClient("http://api_gateway:8080")
	orderUseCase := usecase.NewOrderUsecaseImpl(orderRepository, clientRepo)
	//orderHandler := handler.NewOrderHttpHandler(orderUseCase)

	grpcServer := grpc.NewServer()

	reflection.Register(grpcServer)

	orders.RegisterOrderServiceServer(grpcServer, service.NewOrdersService(orderUseCase))

	return &grpcServerObject{
		server: grpcServer,
		cfg:    conf,
		db:     db,
		log:    log,
	}
}

func (s *grpcServerObject) Start() {
	port := ":50052"
	if s.cfg.Server.Port != "" {
		port = ":" + s.cfg.Server.Port
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		s.log.Fatalf("Failed to listen on %s: %v", port, err)
	}

	s.log.Printf("Starting orders gRPC server on %s", port)
	if err := s.server.Serve(lis); err != nil {
		s.log.Fatalf("Failed to serve: %v", err)
	}
}
