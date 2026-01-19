package main

import (
	"log/slog"
	"os"
	"time"

	"LM-Gate/internal/infra"
	"LM-Gate/internal/logic"
)

func main() {
	// ==================================================
	// 1. إعداد slog (بصيغة نصية للـ Terminal)
	// ==================================================
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// ==================================================
	// 2. Redis Connection (Health Check)
	// ==================================================
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}

	redis := infra.NewRedisService(redisAddr)
	if err := redis.Ping(); err != nil {
		logger.Error("❌ Failed to connect to Redis", slog.Any("error", err))
		os.Exit(1)
	}
	logger.Info("✅ Connected to Redis")

	// ==================================================
	// 3. RabbitMQ Connection (Retry)
	// ==================================================
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

	// ==================================================
	// 6. RabbitMQ Consumers
	// ==================================================

	go logic.RunAPIServer()

	logger.Info("🚀 Server is running and waiting for messages...")
	select {}
}
