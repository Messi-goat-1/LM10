package main

import (
	"log"
	"os"
	"time"

	lmgate "LM-Gate"
	"LM-Gate/handlers" // تأكد من استيراد الـ handlers
	"LM-Gate/services" // تأكد من استيراد الـ services
)

func main() {

	// ==================================================
	// Redis Connection (Health Check)
	// ==================================================
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}

	redis := lmgate.NewRedisService(redisAddr)
	if err := redis.Ping(); err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}

	log.Println("✅ Connected to Redis")

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
			log.Println("✅ Connected to RabbitMQ")
			break
		}

		log.Printf("⏳ RabbitMQ not ready (attempt %d/20): %v", i, err)
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		log.Fatal("❌ Failed to connect to RabbitMQ after multiple attempts")
	}

	defer rabbit.Close()

	// ==================================================
	// Services (تم التحديث)
	// ==================================================
	fileService := services.NewFileService()
	manager := services.NewManager(fileService)

	// ==================================================
	// Event Handlers (تم التحديث للإبقاء على detected فقط)
	// ==================================================
	// نقوم بتعريف الـ Handler الخاص بالملفات المكتشفة فقط
	fileHandler := handlers.NewFileDetectedHandler(manager)

	// ==================================================
	// RabbitMQ Consumer
	// ==================================================
	rabbit.ConsumeMessages("file_events_queue", func(data []byte) {

		// توجيه البيانات مباشرة إلى الـ Handler المخصص لـ Detected
		if err := fileHandler.Handle(data); err != nil {
			log.Printf("❌ Error handling detected file: %v", err)
		}

	})

	log.Println("🚀 Server is running and waiting for messages...")
	select {}
}
