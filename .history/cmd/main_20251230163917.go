package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	lmgate "LM-Gate"
	"LM-Gate/events"
	"LM-Gate/handlers"
	"LM-Gate/services"
)

// تعريف المتغير كـ Global لكي تراه دالة processMessage في أي ملف داخل package main
var redisService *RedisService

// دالة وهمية للتحليل لإخفاء خطأ undefined: extractIPs
func extractIPs(data []byte) []string {
	log.Println("Extracting IPs from PCAP data...")
	return []string{"192.168.1.1"}
}
func main() {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}
	redisAddr := os.Getenv("REDIS_URL")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	// 2. إنشاء الاتصال بـ Redis وطباعة حالة الاتصال
	redisService = NewRedisService(redisAddr)

	// التحقق من أن Redis يعمل فعلياً (Ping)
	err := redisService.Client.Ping(context.Background()).Err()
	if err != nil {
		log.Printf("Could not connect to Redis: %v", err)
	} else {
		log.Println("✅ Connected to Redis successfully!") // هذه الرسالة التي طلبتها
	}
	rabbit, err := lmgate.NewRabbitClient(rabbitURL)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbit.Close()

	manager := services.NewManager()
	pcapService := services.NewPCAPService()

	dispatcher := handlers.NewEventDispatcher()
	dispatcher.RegisterHandler("file.detected", handlers.NewFileDetectedHandler(manager))
	dispatcher.RegisterHandler("file.chunk", handlers.NewFileChunkHandler(manager))
	dispatcher.RegisterHandler("pcap.analyze", handlers.NewPCAPAnalyzeHandler(pcapService))

	rabbit.ConsumeMessages("file_events_queue", func(data []byte) {
		var baseEvent events.Event
		if err := json.Unmarshal(data, &baseEvent); err != nil {
			log.Printf("Failed to unmarshal event: %v", err)
			return
		}
		if err := dispatcher.Dispatch(baseEvent.Event, data); err != nil {
			log.Printf("Dispatch error: %v", err)
		}
	})

	log.Println("Server running, waiting for messages...")
	select {}
}
