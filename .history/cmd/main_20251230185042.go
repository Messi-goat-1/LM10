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

// RabbitSender: هيكل وسيط لربط دالة الرفع بـ RabbitMQ
type RabbitSender struct {
	client *lmgate.RabbitClient
}

// Send: تقوم بإرسال قطعة الملف كحدث عبر RabbitMQ
func (s *RabbitSender) Send(msg lmgate.ChunkMessage) error {
	eventBody := map[string]interface{}{
		"event": "file.chunk",
		"data":  msg,
	}
	body, _ := json.Marshal(eventBody)
	// تصحيح الخطأ: استخدام PublishMessage بدلاً من Publish
	return s.client.PublishMessage("file_events_queue", string(body))
}

func main() {
	// --- إعداد Redis ---
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}
	redis := lmgate.NewRedisService(redisAddr)
	if err := redis.Ping(); err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}
	log.Println("✅ Connected to Redis")

	// --- إعداد RabbitMQ ---
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
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		log.Fatal("❌ Failed to connect to RabbitMQ")
	}
	defer rabbit.Close()

	// --- إعداد الخدمات و Handlers ---
	manager := services.NewManager()
	pcapService := services.NewPCAPService()
	dispatcher := handlers.NewEventDispatcher()

	// ملاحظة: تأكد من تمرير خدمة redis للـ Handler إذا كان يحتاجها
	dispatcher.RegisterHandler("file.chunk", handlers.NewFileChunkHandler(manager))
	dispatcher.RegisterHandler("pcap.analyze", handlers.NewPCAPAnalyzeHandler(pcapService))

	// --- تشغيل مستلم الرسائل ---
	go rabbit.ConsumeMessages("file_events_queue", func(data []byte) {
		var baseEvent events.Event
		if err := json.Unmarshal(data, &baseEvent); err != nil {
			return
		}
		dispatcher.Dispatch(baseEvent.Event, data)
	})

	// --- محاكاة إرسال الملف من سطح المكتب ---
	go func() {
		time.Sleep(10 * time.Second) // انتظار استقرار الحاويات

		// المسار داخل الحاوية (بعد ربط الـ Volumes)
		filePath := "/data/messi.pcap"
		fmt.Printf("🚀 Starting upload for: %s\n", filePath)

		sender := &RabbitSender{client: rabbit}
		chunkSize := int64(512 * 1024) // 512KB

		// استخدام دالة الرفع الحقيقية من مشروعك
		sent, err := lmgate.UploadFile(filePath, chunkSize, sender)
		if err != nil {
			log.Printf("⚠️ Upload error: %v", err)
		} else {
			fmt.Printf("✅ Upload complete! Total chunks: %d\n", sent)
		}
	}()

	log.Println("🚀 Server is running...")
	select {}
}
