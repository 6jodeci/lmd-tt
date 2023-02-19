package main

import (
	"context"
	"lamoda-test/internal/app"
	"lamoda-test/internal/config"
	"lamoda-test/pkg/logging"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger := logging.GetLogger(ctx)

	logger.Info("config initializing")
	cfg := config.GetConfig()

	log.Print("logger initializing")
	ctx = logging.ContextWithLogger(ctx, logger)

	a, err := app.NewApp(ctx, cfg)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("Running Application")
	a.Run(ctx)
}