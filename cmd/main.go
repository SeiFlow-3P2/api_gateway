package main

import (
	"context"
	"log"

	"github.com/SeiFlow-3P2/api_gateway/internal/app"
	"github.com/SeiFlow-3P2/api_gateway/pkg/config"
	"github.com/SeiFlow-3P2/api_gateway/pkg/env"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := env.LoadEnv(); err != nil {
		log.Fatalf("failed to load env: %v", err)
	}

	if env.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	conf, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx := context.Background()
	app := app.NewApp(conf)
	if err := app.Start(ctx); err != nil {
		log.Fatalf("failed to start app: %v", err)
	}
}
