package main

import (
	"context"
	"log"

	"github.com/SeiFlow-3P2/api_gateway/internal/app"
	"github.com/SeiFlow-3P2/api_gateway/internal/config"
)

func main() {
	ctx := context.Background()

	conf, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	app := app.NewApp(conf)
	if err := app.Start(ctx); err != nil {
		log.Fatalf("failed to start app: %v", err)
	}
}
