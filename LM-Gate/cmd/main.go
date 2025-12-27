package main

import (
	lmgate "LM-Gate"
	"LM-Gate/events"
	"LM-Gate/handlers"
	"LM-Gate/services"
	"encoding/json"
	"log"
)

func main() {
	// 1. تهيئة الخدمات (Services) والـ Handlers
	fileService := &services.FileService{}
	fileHandler := handlers.NewFileDetectedHandler(fileService)

	// 2. الاتصال بـ RabbitMQ
	rabbit, _ := lmgate.NewRabbitClient("amqp://guest:guest@localhost:5672/")
	defer rabbit.Close()

	// 3. تعريف كيف سنعالج الرسالة عند وصولها
	messageProcessor := func(data []byte) {
		var event events.FileDetectedEvent

		// تحويل JSON إلى Struct
		if err := json.Unmarshal(data, &event); err != nil {
			log.Printf("خطأ في تحليل البيانات: %v", err)
			return
		}

		// استدعاء الـ Handler الذي قمت بكتابته
		fileHandler.Handle(event)
	}

	// 4. البدء بالاستماع
	rabbit.ConsumeMessages("file_events_queue", messageProcessor)

	// منع البرنامج من الإغلاق
	select {}
}

/*
msg := <-rabbitMessages

err := DispatchEvent(msg.Type, msg.Body)
if err != nil {
    log.Println("event error:", err)
}
*/
