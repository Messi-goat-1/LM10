package main

/*
func main() {
	// 1. تهيئة الخدمات والمعالجات (استخدام الاسم الجديد Manager)
	fileManager := services.NewManager() // تم التعريف باسم fileManager

	// تمرير fileManager لكل المعالجات لضمان توافق الأنواع
	fileHandler := handlers.NewFileDetectedHandler(fileManager)
	chunkHandler := handlers.NewFileChunkHandler(fileManager)

	// إنشاء الموزع (Dispatcher) لترتيب الكود
	dispatcher := handlers.NewEventDispatcher(fileHandler, chunkHandler)

	// 2. الاتصال بـ RabbitMQ
	rabbit, _ := lmgate.NewRabbitClient("amqp://guest:guest@localhost:5672/")
	defer rabbit.Close()

	// 3. معالج رسائل موحد وبسيط (Message Processor)
	messageProcessor := func(data []byte) {
		// ملاحظة: الـ routingKey يفضل جلبه من خصائص الرسالة في RabbitMQ
		// سنفترض حالياً "file.detected" للتجربة
		routingKey := "file.detected"

		err := dispatcher.Dispatch(routingKey, data)
		if err != nil {
			log.Printf("❌ Dispatch error: %v", err)
		}
	}

	// 4. الاستماع (Consume) للأحداث من الطابور
	// تأكد أن اسم الطابور "file_events_queue" مطابق لما في RabbitMQ
	rabbit.ConsumeMessages("file_events_queue", messageProcessor)

	log.Println("🚀 السيرفر يعمل الآن ويستمع للأحداث...")
	select {}
}
*/
