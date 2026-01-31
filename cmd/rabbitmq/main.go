package main

import (
	"LM-Gate/internal/infra"
	"log/slog"
	"os"
	"time"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@rabbitmq:5672/"
	}

	var rabbit *infra.RabbitClient
	var err error
	for i := 1; i <= 20; i++ {
		rabbit, err = infra.NewRabbitClient(rabbitURL)
		if err == nil {
			logger.Info("✅ Connected to RabbitMQ")
			break
		}
		logger.Warn("⏳ RabbitMQ not ready", slog.Int("attempt", i), slog.Any("error", err))
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		logger.Error("❌ Failed to connect to RabbitMQ after multiple attempts")
		os.Exit(1)
	}
	defer rabbit.Close()

	logger.Info("🚀 RabbitMQ Worker is running and waiting...")
	select {}
}
