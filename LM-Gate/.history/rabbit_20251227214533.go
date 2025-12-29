package lmgate

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitClient هيكل لتخزين بيانات الاتصال
type RabbitClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewRabbitClient دالة لإنشاء اتصال جديد وتجهيزه
func NewRabbitClient(url string) (*RabbitClient, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitClient{
		conn:    conn,
		channel: ch,
	}, nil
}

// Close لإغلاق الاتصال بشكل آمن
func (r *RabbitClient) Close() {
	r.channel.Close()
	r.conn.Close()
}

// PublishMessage دالة مخصصة لإرسال الرسائل
func (r *RabbitClient) PublishMessage(queueName string, message string) error {
	_, err := r.channel.QueueDeclare(
		queueName, false, false, false, false, nil,
	)
	if err != nil {
		return err
	}

	return r.channel.PublishWithContext(context.Background(),
		"", queueName, false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

// دالة منفصلة لإدارة التكرار (Logic)
func RunHeartbeat(rabbit *RabbitClient) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		msg := fmt.Sprintf("Server is alive at %s", time.Now().Format("15:04:05"))

		if err := rabbit.PublishMessage("server_status", msg); err != nil {
			log.Printf("خطأ في الإرسال: %v", err)
		} else {
			fmt.Println("Sent: Heartbeat to RabbitMQ")
		}
	}
}

// ConsumeMessages دالة للاستماع للرسائل من طابور معين
// ConsumeMessages تستقبل الآن دالة لمعالجة البيانات المستلمة
func (r *RabbitClient) ConsumeMessages(queueName string, processor func([]byte)) {
	msgs, err := r.channel.Consume(
		queueName, "", true, false, false, false, nil,
	)
	if err != nil {
		log.Printf("فشل بدء الاستماع: %v", err)
		return
	}

	go func() {
		for d := range msgs {
			// نرسل البيانات الخام للدالة المعالجة
			processor(d.Body)
		}
	}()
}
