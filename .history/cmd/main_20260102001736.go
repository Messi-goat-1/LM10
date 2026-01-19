package main

import (
	"log/slog" // إضافة مكتبة slog
	"os"
	"time"

	lmgate "LM-Gate"
	"LM-Gate/handlers"
	"LM-Gate/services"
)

func main() {
	// ==================================================
	// إعداد slog (بصيغة نصية للـ Terminal)
	// ==================================================
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// ==================================================
	// Redis Connection (Health Check)
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
	// RabbitMQ Connection (Retry)
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
	// Services (تمرير الـ logger للخدمة)
	// ==================================================
	fileService := services.NewFileService(logger)
	manager := services.NewManager(fileService)

	// ==================================================
	// Event Handlers
	// ==================================================
	fileHandler := handlers.NewFileDetectedHandler(manager)

	// ==================================================
	// RabbitMQ Consumer
	// ==================================================
	rabbit.ConsumeMessages("file_events_queue", func(data []byte) {
		// توجيه البيانات مباشرة إلى الـ Handler المخصص لـ Detected
		if err := fileHandler.Handle(data); err != nil {
			logger.Error("❌ Error handling detected file", slog.Any("error", err))
		}
	})

	logger.Info("🚀 Server is running and waiting for messages...")
	select {}
}
