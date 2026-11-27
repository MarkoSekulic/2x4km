package main

import (
	"context"
	"example/gateway/config"
	"example/gateway/middleware"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"proto/notifications-service"
	"proto/template-service"
	"proto/users-service"
	"syscall"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func dialConnection(addr string) *grpc.ClientConn {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	return conn
}

func main() {

	fmt.Println("Starting api gateway...")

	cfg := config.GetConfig()

	ctx := context.Background()

	// Apply Response Middleware
	gwmux := runtime.NewServeMux(runtime.WithMetadata(middleware.CookieToMetadataMiddleware), runtime.WithForwardResponseOption(middleware.ResponseHeaderMiddleware))

	// Register Template
	templateConn := dialConnection(cfg.TemplateServiceAddress)
	defer templateConn.Close()

	err := template.RegisterTemplateServiceHandlerClient(ctx, gwmux, template.NewTemplateServiceClient(templateConn))
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}
	// Register User
	usersConn := dialConnection(cfg.UsersServiceAddress)
	defer usersConn.Close()

	err = users.RegisterUsersServiceHandlerClient(ctx, gwmux, users.NewUsersServiceClient(usersConn))
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}
	// Register Notification
	notificationsConn := dialConnection(cfg.NotificationsServiceAddress)
	defer usersConn.Close()

	err = notifications.RegisterNotificationsServiceHandlerClient(ctx, gwmux, notifications.NewNotificationsServiceClient(notificationsConn))
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}
	// Register Blank
	// ...

	// Apply Middleware

	gwhandler := middleware.RegisterLoggerMiddleware(gwmux)
	// gwhandler = middleware.RegisterLoggerMiddleware(gwhandler)
	// gwhandler = middleware.RegisterLoggerMiddleware(gwhandler)
	// ...
	gwhandler = middleware.CorsMiddleware(gwhandler)

	// Rate Limiter (10 per 10 second)
	lmt := tollbooth.NewLimiter(1, &limiter.ExpirableOptions{
		DefaultExpirationTTL: time.Hour,
	}).
		SetBurst(10).
		SetIPLookups([]string{"X-Real-IP", "X-Forwarded-For", "RemoteAddr"}).
		SetMessageContentType("application/json").
		SetMessage(`{"code": 8, "message": "rate limit exceeded"}`)
	gwhandler = tollbooth.LimitHandler(lmt, gwhandler)

	gwServer := &http.Server{
		Addr:    cfg.Address,
		Handler: gwhandler,
	}

	go func() {
		if err := gwServer.ListenAndServeTLS("../certs/localhost.pem", "../certs/localhost-key.pem"); err != nil {
			log.Fatal("server error: ", err)
		}
	}()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGTERM)

	<-stopCh

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	if err = gwServer.Shutdown(ctx); err != nil {
		log.Fatalln("error while stopping server: ", err)
	}
	defer cancel()
}
