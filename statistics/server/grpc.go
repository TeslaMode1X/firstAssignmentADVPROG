package server

import (
	"context"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/statistics/internal/model"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/statistics/pkg/nats/producer"
	"log"
	"net"
	"time"

	"github.com/TeslaMode1X/firstAssignmentADVPROG/proto/gen/statistics"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/statistics/config"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/statistics/internal/interfaces"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/statistics/internal/repository"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/statistics/internal/service"
	statGrpc "github.com/TeslaMode1X/firstAssignmentADVPROG/statistics/internal/service/grpc"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/statistics/pkg/nats"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/statistics/pkg/nats/consumer"
	natsHandler "github.com/TeslaMode1X/firstAssignmentADVPROG/statistics/pkg/nats/handler"

	"google.golang.org/grpc"
)

type grpcServerObject struct {
	server            *grpc.Server
	cfg               *config.Config
	db                interfaces.Database
	log               *log.Logger
	natsClient        *nats.Client
	pubSub            *consumer.PubSub
	sendRefreshCancel context.CancelFunc
}

func NewGrpcServerObject(conf *config.Config, db interfaces.Database, log *log.Logger) Server {
	statisticsRepo := repository.NewStatisticsRepo(db)
	statisticsService := service.NewStatisticsService(statisticsRepo)

	grpcServer := grpc.NewServer()

	s := &grpcServerObject{
		server: grpcServer,
		cfg:    conf,
		db:     db,
		log:    log,
	}

	var err error
	s.natsClient, err = nats.NewClient(context.Background(), []string{"nats_server:4222"}, "", true)
	if err != nil {
		log.Fatalf("failed to connect to NATS: %v", err)
	}
	log.Println("NATS connection status is", s.natsClient.Conn.Status().String())

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

	orderProducer := producer.NewOrderProducer(s.natsClient)

	s.startCacheRefreshJob(context.Background(), *orderProducer, 1*time.Minute)

	statistics.RegisterStatisticsServiceServer(grpcServer, statGrpc.NewStatisticsService(statisticsService))

	return s
}

func (s *grpcServerObject) Start() {
	port := ":50053"
	if s.cfg.Server.Port != "" {
		port = ":" + s.cfg.Server.Port
	}

	errCh := make(chan error, 1)
	s.pubSub.Start(context.Background(), errCh)

	go func() {
		for err := range errCh {
			s.log.Printf("NATS error: %v", err)
		}
	}()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		s.log.Fatalf("Failed to listen on %s: %v", port, err)
	}

	s.log.Printf("Starting statistics gRPC server on %s", port)
	if err := s.server.Serve(lis); err != nil {
		s.log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *grpcServerObject) Stop() {
	if s.natsClient != nil {
		s.log.Println("Closing NATS connection...")
		s.natsClient.CloseConnect()
	}
	if s.server != nil {
		s.server.GracefulStop()
	}
}

func (s *grpcServerObject) startCacheRefreshJob(ctx context.Context, orderProducer producer.OrderProducer, refreshInterval time.Duration) {
	refreshCtx, cancel := context.WithCancel(ctx)
	s.sendRefreshCancel = cancel

	pr := model.Product{}

	go func() {
		ticker := time.NewTicker(refreshInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.log.Println("Running scheduled nats sending...")
				err := orderProducer.Push(ctx, pr)
				if err != nil {
					s.log.Printf("Scheduled nats refresh failed: %v", err)
				} else {
					s.log.Println("Scheduled nats refresh completed successfully")
				}
			case <-refreshCtx.Done():
				s.log.Println("Cache refresh job terminated")
				return
			}
		}
	}()

	s.log.Printf("Background cache refresh job started with %v interval", refreshInterval)
}
