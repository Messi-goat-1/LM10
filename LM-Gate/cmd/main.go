package main

import (
	"log/slog"
	"os"
	"time"

	lmgate "LM-Gate"
	api "LM-Gate/API"
	"LM-Gate/handlers"
	"LM-Gate/services"
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

	redis := lmgate.NewRedisService(redisAddr)
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

	var rabbit *lmgate.RabbitClient
	var err error
	for i := 1; i <= 20; i++ {
		rabbit, err = lmgate.NewRabbitClient(rabbitURL)
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
	// 4. Services (حقن الـ logger والخدمات)
	// ==================================================
	fileService := services.NewFileService(logger)
	manager := services.NewManager(fileService)

	// ==================================================
	// 5. Event Handlers
	// ==================================================
	fileDetectedHandler := handlers.NewFileDetectedHandler(manager)
	fileCollectionHandler := handlers.NewFileCollectionHandler(manager)

	// ==================================================
	// 6. RabbitMQ Consumers
	// ==================================================

	// استهلاك رسائل اكتشاف الملفات
	rabbit.ConsumeMessages("file_events_queue", func(data []byte) {
		if err := fileDetectedHandler.Handle(data); err != nil {
			logger.Error("❌ Error handling detected file", slog.Any("error", err))
		}
	})

	// استهلاك رسائل تجميع الملفات
	rabbit.ConsumeMessages("file_collection_queue", func(data []byte) {
		if err := fileCollectionHandler.Handle(data); err != nil {
			logger.Error("❌ Error handling collection file", slog.Any("error", err))
		}
	})
	go api.RunAPIServer()

	logger.Info("🚀 Server is running and waiting for messages...")
	select {}
}
