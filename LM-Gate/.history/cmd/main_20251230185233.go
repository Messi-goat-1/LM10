package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	lmgate "LM-Gate"
	"LM-Gate/events"
	"LM-Gate/handlers"
	"LM-Gate/services"
)

// RabbitSender يربط دالة الرفع بـ RabbitMQ باستخدام الدوال المتوفرة لديك
type RabbitSender struct {
	client *lmgate.RabbitClient
}

func (s *RabbitSender) Send(msg lmgate.ChunkMessage) error {
	eventBody := map[string]interface{}{
		"event": "file.chunk",
		"data":  msg,
	}
	body, _ := json.Marshal(eventBody)
	// استخدام PublishMessage المعرفة في ملف rabbit.go الخاص بك [cite: 2]
	return s.client.PublishMessage("file_events_queue", string(body))
}

func main() {
	// 1. إعداد Redis [cite: 1]
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}
	redis := lmgate.NewRedisService(redisAddr)
	if err := redis.Ping(); err != nil {
		log.Fatalf("❌ Redis Error: %v", err)
	}

	// 2. إعداد RabbitMQ [cite: 1, 2]
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@rabbitmq:5672/"
	}
	var rabbit *lmgate.RabbitClient
	var err error
	for i := 1; i <= 10; i++ {
		rabbit, err = lmgate.NewRabbitClient(rabbitURL)
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("❌ RabbitMQ Error")
	}
	defer rabbit.Close()

	// 3. تهيئة الموزع والخدمات [cite: 1, 3]
	manager := services.NewManager()
	pcapService := services.NewPCAPService()
	dispatcher := handlers.NewEventDispatcher()

	// ربط الـ Handler مع تمرير خدمة الـ Redis المطلوبة للتخزين
	dispatcher.RegisterHandler("file.chunk", handlers.NewFileChunkHandler(manager))
	dispatcher.RegisterHandler("pcap.analyze", handlers.NewPCAPAnalyzeHandler(pcapService))

	// 4. تشغيل المستلم (Consumer) [cite: 2]
	go rabbit.ConsumeMessages("file_events_queue", func(data []byte) {
		var baseEvent events.Event
		json.Unmarshal(data, &baseEvent)
		dispatcher.Dispatch(baseEvent.Event, data)
	})

	// 5. تشغيل دورة الأحداث آلياً
	go func() {
		time.Sleep(10 * time.Second) // انتظار استقرار الحاويات

		// المسار الذي حددناه في docker-compose
		filePath := "/data/messi.pcap"
		fmt.Printf("🚀 بدء رفع الملف: %s\n", filePath)

		sender := &RabbitSender{client: rabbit}
		// استخدام دالة الرفع الحقيقية من ملف client.go
		sent, err := lmgate.UploadFile(filePath, 512*1024, sender)
		if err != nil {
			log.Printf("⚠️ فشل الرفع: %v", err)
		} else {
			fmt.Printf("✅ تم إرسال %d قطعة وبدأت دورة التحليل!\n", sent)
		}
	}()

	log.Println("🚀 السيرفر يعمل...")
	select {}
}
