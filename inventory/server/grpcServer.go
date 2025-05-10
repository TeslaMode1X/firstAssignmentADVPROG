package server

import (
	"context"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/config"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/database"
	redisCache "github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/database/cache"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/repository"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/service"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/usecase"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/pkg/nats"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/pkg/nats/producer"
	redisconn "github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/pkg/redis"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/proto/gen/inventory"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/proto/gen/promotion"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"
)

type grpcServerObject struct {
	server             *grpc.Server
	cfg                *config.Config
	db                 database.Database
	log                *log.Logger
	natsClient         *nats.Client
	cacheRefreshCancel context.CancelFunc
}

func NewGRPCServer(conf *config.Config, db database.Database, log *log.Logger) Server {
	productRepository := repository.NewProductPostgresRepository(db)
	promoteRepository := repository.NewPromotePostgresRepository(db)

	ctx := context.Background()

	// Create NATS client
	natsClient, err := nats.NewClient(context.Background(), []string{"nats_server:4222"}, "", true) // Remove the NKey if not needed
	if err != nil {
		log.Fatal(err)
	}
	log.Println("NATS connection status is", natsClient.Conn.Status().String())

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

	inventoryProducer := producer.NewInventoryProducer(natsClient)

	productUseCase := usecase.NewProductUsecaseImpl(productRepository, promoteRepository, clientRedisCache)
	promotionUseCase := usecase.NewPromoteUsecaseImpl(productRepository, promoteRepository)

	log.Println("Initializing Redis cache with all products...")
	err = productUseCase.RefreshCache(ctx)
	if err != nil {
		log.Printf("Warning: Failed to initialize product cache: %v", err)
	} else {
		log.Println("Product cache successfully initialized")
	}

	grpcServer := grpc.NewServer()

	reflection.Register(grpcServer)

	inventory.RegisterInventoryServiceServer(grpcServer, service.NewInventoryService(productUseCase, inventoryProducer))
	promotion.RegisterPromotionServiceServer(grpcServer, service.NewPromotionService(promotionUseCase))

	server := &grpcServerObject{
		server:     grpcServer,
		cfg:        conf,
		db:         db,
		log:        log,
		natsClient: natsClient,
	}

	server.startCacheRefreshJob(ctx, productUseCase, 12*time.Hour)

	return server
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

func (s *grpcServerObject) Stop() {
	if s.cacheRefreshCancel != nil {
		s.cacheRefreshCancel()
	}
	s.server.GracefulStop()
}

func (s *grpcServerObject) startCacheRefreshJob(ctx context.Context, productUseCase *usecase.ProductUsecaseImpl, refreshInterval time.Duration) {
	refreshCtx, cancel := context.WithCancel(ctx)
	s.cacheRefreshCancel = cancel

	go func() {
		ticker := time.NewTicker(refreshInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.log.Println("Running scheduled cache refresh...")
				err := productUseCase.RefreshCache(refreshCtx)
				if err != nil {
					s.log.Printf("Scheduled cache refresh failed: %v", err)
				} else {
					s.log.Println("Scheduled cache refresh completed successfully")
				}
			case <-refreshCtx.Done():
				s.log.Println("Cache refresh job terminated")
				return
			}
		}
	}()

	s.log.Printf("Background cache refresh job started with %v interval", refreshInterval)
}
