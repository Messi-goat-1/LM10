package lmgate

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitClient represents a RabbitMQ client.
// It holds the connection and channel used for messaging.
//
// NOTE: This struct is responsible only for messaging logic.
// TODO: Add automatic reconnection support.
type RabbitClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewRabbitClient creates a new RabbitMQ connection
// and prepares a channel for publishing and consuming messages.
//
// NOTE: Call Close() when the client is no longer needed.
// FIXME: No custom timeout is configured for the connection.
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

// Close safely closes the RabbitMQ channel and connection.
//
// NOTE: Should be called when shutting down the server.
// TODO: Add nil checks before closing.
func (r *RabbitClient) Close() {
	r.channel.Close()
	r.conn.Close()
}

// PublishMessage sends a text message to a specific queue.
// PublishMessage sends a text message to a specific queue.
//
// NOTE: The queue is created if it does not already exist.
// Now updated to use durable = true to match ConsumeMessages.
func (r *RabbitClient) PublishMessage(queueName string, message string) error {
	// 1. التأكد من وجود الطابور بنفس إعدادات المستهلك (Durable: true)
	_, err := r.channel.QueueDeclare(
		queueName,
		true,  // durable: يجب أن تكون true لتتطابق مع دالة الاستقبال
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	// 2. إرسال الرسالة إلى الطابور
	return r.channel.PublishWithContext(
		context.Background(),
		"",        // exchange
		queueName, // routing key
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
			// نصيحة إضافية: تجعل الرسالة نفسها محفوظة على القرص
			DeliveryMode: amqp.Persistent,
		},
	)
}

// RunHeartbeat sends a heartbeat message every minute.
//
// NOTE: Used to signal that the server is alive.
// TODO: Make the interval configurable.
// FIXME: Stops only when the application exits.
func RunHeartbeat(rabbit *RabbitClient) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		msg := fmt.Sprintf(
			"Server is alive at %s",
			time.Now().Format("15:04:05"),
		)

		if err := rabbit.PublishMessage("server_status", msg); err != nil {
			log.Printf("Failed to send heartbeat: %v", err)
		} else {
			fmt.Println("Sent: Heartbeat to RabbitMQ")
		}
	}
}

// ConsumeMessages listens to messages from a specific queue
// and passes the received data to a processor function.
//
// NOTE: Now it ensures the queue exists before consuming (Durable = true).
func (r *RabbitClient) ConsumeMessages(
	queueName string,
	processor func([]byte),
) {
	// 1. التأكد من وجود الطابور وتفعيل خاصية البقاء (Durable: true)
	// هذا يمنع خطأ "queue not found" ويحمي الطابور من الحذف عند رسترت الأرنب
	_, err := r.channel.QueueDeclare(
		queueName,
		true,  // durable: الطابور يبقى محفوظاً في القرص
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		log.Printf("Failed to declare queue: %v", err)
		return
	}

	// 2. البدء باستهلاك الرسائل من الطابور
	msgs, err := r.channel.Consume(
		queueName,
		"",    // consumer tag
		true,  // auto-ack: تمكين التأكيد التلقائي لاستلام الرسالة
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		log.Printf("Failed to start consuming: %v", err)
		return
	}

	// 3. تشغيل مستمع الرسائل في خلفية البرنامج (Goroutine)
	go func() {
		log.Printf("Successfully started consuming from: %s", queueName)
		for d := range msgs {
			// تمرير بيانات الرسالة الخام إلى دالة المعالجة
			processor(d.Body)
		}
	}()
}

// في ملف rabbit.go
func (r *RabbitClient) Publish(queueName string, body []byte) error {
	return r.Channel.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
