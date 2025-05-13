package server

import (
	"context"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/config"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/internal/database"
	redisCache "github.com/TeslaMode1X/firstAssignmentADVPROG/orders/internal/database/cache"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/internal/handler/client"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/internal/repository"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/internal/service"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/internal/usecase"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/pkg/nats"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/pkg/nats/producer"
	redisconn "github.com/TeslaMode1X/firstAssignmentADVPROG/orders/pkg/redis"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/proto/gen/orders"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"
)

type grpcServerObject struct {
	server     *grpc.Server
	cfg        *config.Config
	db         database.Database
	log        *log.Logger
	natsClient *nats.Client
}

func NewGRPCServer(conf *config.Config, db database.Database, log *log.Logger) Server {
	orderRepository := repository.NewOrderPostgresRepository(db)
	clientRepo := client.NewInventoryClient("http://api_gateway:8080")

	ctx := context.Background()

	// Create NATS client
	natsClient, err := nats.NewClient(context.Background(), []string{"nats_server:4222"}, "", true) // Remove the NKey if not needed
	if err != nil {
		log.Fatal(err)
	}
	log.Println("NATS connection status is", natsClient.Conn.Status().String())

	OrdersProducer := producer.NewOrdersProducer(natsClient)

	grpcServer := grpc.NewServer()

	reflection.Register(grpcServer)

	// REDIS connection
	log.Println("Attempting to connect to Redis...")
	redisClient, err := redisconn.NewClient(ctx, redisconn.GetRedisConfig())
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Successfully connected to Redis!")

	// Test Redis connection with PING
	pingErr := redisClient.Ping(ctx)
	if pingErr != nil {
		log.Fatalf("Redis PING failed: %v", pingErr)
	}
	log.Println("Redis PING successful - connection is working!")

	// REDIS cache
	clientRedisCache := redisCache.NewClient(redisClient, 12*time.Hour)
	log.Println("Redis cache client initialized with 10 hour TTL")

	orderUseCase := usecase.NewOrderUsecaseImpl(orderRepository, clientRepo, clientRedisCache)

	orders.RegisterOrderServiceServer(grpcServer, service.NewOrdersService(orderUseCase, OrdersProducer))

	// REMOVE the defer natsClient.Close() line!

	return &grpcServerObject{
		server:     grpcServer,
		cfg:        conf,
		db:         db,
		log:        log,
		natsClient: natsClient, // Store the NATS client
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

func (s *grpcServerObject) Stop() {
	if s.natsClient != nil {
		s.log.Println("Closing NATS connection...")
		s.natsClient.Close()
	}
	s.server.GracefulStop()
}
