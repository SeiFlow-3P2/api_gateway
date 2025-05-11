package main

import (
	"log"
	"os"

	"github.com/SeiFlow-3P2/api_gateway/internal/app"
	"github.com/SeiFlow-3P2/api_gateway/internal/config"
)

func main() {
	log := log.New(os.Stdout, "", log.LstdFlags)

	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	a := app.NewApp(cfg)
	if err := a.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
