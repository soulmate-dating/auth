package main

import (
	"context"
	"github.com/soulmate-dating/auth/internal/adapters/jwt"
	"github.com/soulmate-dating/auth/internal/adapters/postgres"
	"github.com/soulmate-dating/auth/internal/app"
	"github.com/soulmate-dating/auth/internal/graceful"
	grpcSvc "github.com/soulmate-dating/auth/internal/ports/grpc"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"os"

	"log"
	"net"
)

const (
	grpcPort = ":8081"
)

func main() {
	ctx := context.Background()

	dbConn, err := postgres.Connect()
	jwt := jwt.NewWrapper("secret-key", "auth-service", 5, 30*24*60)
	if err != nil {
		log.Fatal(err.Error())
	}
	appSvc := app.NewApp(dbConn, jwt)

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	svc := grpcSvc.NewService(appSvc)
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		grpcSvc.UnaryLoggerInterceptor,
		grpcSvc.UnaryRecoveryInterceptor(),
	))
	grpcSvc.RegisterAuthServiceServer(grpcServer, svc)

	eg, ctx := errgroup.WithContext(ctx)

	sigQuit := make(chan os.Signal, 1)
	eg.Go(graceful.CaptureSignal(ctx, sigQuit))
	// run grpc server
	eg.Go(grpcSvc.RunGRPCServerGracefully(ctx, lis, grpcServer))

	if err := eg.Wait(); err != nil {
		log.Printf("gracefully shutting down the servers: %s\n", err.Error())
	}
	log.Println("servers were successfully shutdown")
}
