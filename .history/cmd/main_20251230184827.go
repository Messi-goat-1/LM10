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

// RabbitSender: يربط دالة الرفع بمحرك RabbitMQ الخاص بك
type RabbitSender struct {
	client *lmgate.RabbitClient
}

// Send: تقوم بتحويل رسالة القطعة إلى حدث وإرسالها
func (s *RabbitSender) Send(msg lmgate.ChunkMessage) error {
	eventBody := map[string]interface{}{
		"event": "file.chunk",
		"data":  msg,
	}
	body, _ := json.Marshal(eventBody)
	// استخدام PublishMessage المعرفة في ملف rabbit.go
	return s.client.PublishMessage("file_events_queue", string(body))
}

func main() {
	// 1. إعداد الاتصال بـ Redis
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}
	redis := lmgate.NewRedisService(redisAddr)
	if err := redis.Ping(); err != nil {
		log.Fatalf("❌ فشل الاتصال بـ Redis: %v", err)
	}
	log.Println("✅ متصل بـ Redis")

	// 2. إعداد الاتصال بـ RabbitMQ
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@rabbitmq:5672/"
	}
	var rabbit *lmgate.RabbitClient
	var err error
	for i := 1; i <= 10; i++ {
		rabbit, err = lmgate.NewRabbitClient(rabbitURL)
		if err == nil {
			log.Println("✅ متصل بـ RabbitMQ")
			break
		}
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("❌ فشل الاتصال بـ RabbitMQ")
	}
	defer rabbit.Close()

	// 3. تهيئة الخدمات والموزع (Dispatcher)
	manager := services.NewManager()
	pcapService := services.NewPCAPService()
	dispatcher := handlers.NewEventDispatcher()

	// تسجيل الـ Handlers وربطهم بخدمة Redis
	dispatcher.RegisterHandler("file.chunk", handlers.NewFileChunkHandler(manager))
	dispatcher.RegisterHandler("pcap.analyze", handlers.NewPCAPAnalyzeHandler(pcapService))

	// 4. تشغيل مستمع الأحداث
	go rabbit.ConsumeMessages("file_events_queue", func(data []byte) {
		var baseEvent events.Event
		if err := json.Unmarshal(data, &baseEvent); err != nil {
			return
		}
		dispatcher.Dispatch(baseEvent.Event, data)
	})

	// 5. محاكاة رفع ملف حقيقي باستخدام دالة UploadFile من مشروعك
	go func() {
		time.Sleep(10 * time.Second) // انتظار استقرار Docker

		fmt.Println("🚀 بدء عملية رفع الملف التجريبي...")
		sender := &RabbitSender{client: rabbit}

		// استخدام دوال client.go الموجودة لديك
		// تأكد من وجود ملف بهذا الاسم في مجلد المشروع
		fileName := "/home/messi/Desktop/messi.pcap"
		chunkSize := int64(512 * 1024) // 512KB للقطعة

		sent, err := lmgate.UploadFile(fileName, chunkSize, sender)
		if err != nil {
			log.Printf("⚠️ فشل الرفع (تأكد من وجود ملف %s): %v", fileName, err)
		} else {
			fmt.Printf("✅ تم إرسال %d قطعة، وتم تفعيل دورة الأحداث بنجاح!\n", sent)
		}
	}()

	log.Println("🚀 السيرفر يعمل وينتظر الأحداث...")
	select {}
}
