package main

import (
	"2x4km/template/config"
	"2x4km/template/handlers"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"proto/template-service"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	fmt.Println("Starting template service...")

	cfg := config.GetConfig()

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		log.Fatalln(err)
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(listener)

	// Bootstrap gRPC server.
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	// Bootstrap gRPC service server and respond to request.
	templateHandler := handlers.TemplateHandler{}
	template.RegisterTemplateServiceServer(grpcServer, templateHandler)

	// Graceful Shutdown
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal("server error: ", err)
		}
	}()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGTERM)

	<-stopCh

	grpcServer.Stop()
}
