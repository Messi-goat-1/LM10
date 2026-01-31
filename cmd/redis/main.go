package main

import (
	"LM-Gate/internal/infra"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

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

	select {}
}
