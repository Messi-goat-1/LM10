package main

import (
	"LM-Gate/analysis" // استبدل project_name باسم مشروعك في go.mod
	"fmt"
	"log"
)

func main() {
	// 1. حدد مسار ملف PCAP حقيقي موجود على جهازك حالياً
	testFilePath := "/home/messi/Desktop/messi"

	fmt.Println("🧪 بدأت تجربة فحص البيانات...")

	// 2. استدعاء دالة الفتح التي فصلناها
	handle, err := analysis.GetFileHandle(testFilePath)
	if err != nil {
		log.Fatalf("❌ فشل الاختبار في مرحلة الفتح: %v", err)
	}
	defer handle.Close()

	// 3. استدعاء دالة التحليل الشاملة لرؤية البيانات
	fmt.Println("-------------------------------------------")
	analysis.RunFullAnalysis(handle)
	fmt.Println("-------------------------------------------")

	fmt.Println("✅ انتهى الاختبار بنجاح.")
}

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
