package main

import (
	"encoding/json"
	"log"
	"os"

	lmgate "LM-Gate"
	"LM-Gate/events"
	"LM-Gate/handlers"
	"LM-Gate/services"
)

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redis := lmgate.NewRedisService(redisAddr)
	if err := redis.Ping(); err != nil {
		log.Fatalf("❌ فشل الاتصال بـ Redis: %v", err)
	}
	log.Println("✅ تم الاتصال بـ Redis")

	// RabbitMQ (موجود كما هو)
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	rabbit, err := lmgate.NewRabbitClient(rabbitURL)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbit.Close()

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
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
