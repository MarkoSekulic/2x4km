package main

import (
	"2x4km/files/config"
	"2x4km/files/handlers"
	"2x4km/files/repo"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"proto/files-service"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {

	fmt.Println("Starting files service...")

	cfg := config.GetConfig()
	ctx := context.Background()

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		log.Fatalln(err)
	}

	// Bootstrap gRPC server.
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	// Init Repository and Handler
	filesRepo, disconnectDB, err := repo.InitDB(ctx, cfg.DBConfig)
	if err != nil {
		log.Fatalln(err)
	}
	defer disconnectDB(ctx)

	filesHandler := handlers.InitFilesHandler(filesRepo, &cfg.HandlerConfig)

	// Bootstrap gRPC service server and respond to request.
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())
	files.RegisterFilesServiceServer(grpcServer, filesHandler)
	// Graceful Shutdown
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal("server error: ", err)
		}
	}()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGTERM, syscall.SIGINT)

	<-stopCh

	grpcServer.GracefulStop()
}
