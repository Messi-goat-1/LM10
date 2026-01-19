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

func main() {
	// 1. الاتصال بـ Redis
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}
	redis := lmgate.NewRedisService(redisAddr)
	if err := redis.Ping(); err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}
	log.Println("✅ Connected to Redis")

	// 2. الاتصال بـ RabbitMQ مع محاولات إعادة الاتصال
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

	// 3. إعداد الخدمات والـ Dispatcher
	manager := services.NewManager()
	pcapService := services.NewPCAPService()
	dispatcher := handlers.NewEventDispatcher()

	// ربط الـ Handlers مع تمرير خدمة Redis
	dispatcher.RegisterHandler("file.chunk", handlers.NewFileChunkHandler(manager))
	dispatcher.RegisterHandler("pcap.analyze", handlers.NewPCAPAnalyzeHandler(pcapService))

	// 4. تشغيل المستلم (Consumer)
	go rabbit.ConsumeMessages("file_events_queue", func(data []byte) {
		var baseEvent events.Event
		json.Unmarshal(data, &baseEvent)
		dispatcher.Dispatch(baseEvent.Event, data)
	})

	// 5. تشغيل المحاكي لإرسال ملف حقيقي للتجربة
	go func() {
		time.Sleep(10 * time.Second) // انتظار استقرار الحاويات
		sendRealFile(rabbit, "messi.pcap")
	}()

	log.Println("🚀 Server is running...")
	select {}
}

// دالة إرسال ملف حقيقي وتقسيمه لأحداث
func sendRealFile(rabbit *lmgate.RabbitClient, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("⚠️ لم يتم العثور على ملف التجربة: %v", err)
		return
	}
	defer file.Close()

	fileID := "test-pcap-001"
	buffer := make([]byte, 512*1024) // تقسيم لقطع بحجم 512 كيلوبايت

	for {
		n, err := file.Read(buffer)
		isEOF := err != nil

		chunkEvent := map[string]interface{}{
			"event": "file.chunk",
			"data": map[string]interface{}{
				"FileID": fileID,
				"Data":   buffer[:n],
				"IsEOF":  isEOF,
			},
		}

		body, _ := json.Marshal(chunkEvent)
		rabbit.PublishBytes("file_events_queue", body)
		fmt.Printf("📦 تم إرسال قطعة: %d bytes (IsEOF: %v)\n", n, isEOF)

		if isEOF {
			break
		}
	}
}
