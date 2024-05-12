package main

import (
	"context"
	"log"

	"github.com/soulmate-dating/auth/internal/app"
	"github.com/soulmate-dating/auth/internal/config"
	"github.com/soulmate-dating/auth/internal/ports/grpc"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	appSvc := app.New(ctx, cfg)
	grpc.Run(ctx, cfg, appSvc)
}
