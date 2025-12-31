package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/barbodimani81/Event-DDD.git/internal/domain"
	"github.com/barbodimani81/Event-DDD.git/internal/infra/rabbit"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// 1. Connect to RabbitMQ (Same as API)
	rabbitURL := os.Getenv("RABBIT_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}
	conn, err := amqp.Dial(rabbitURL)

	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// 2. Initialize the Consumer
	consumer, err := rabbit.NewRabbitConsumer(conn)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}

	// 3. Start Listening
	msgs, err := consumer.Start()
	if err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}

	log.Println(" [*] Worker started. Waiting for messages...")

	// 4. Create a channel to keep the program running forever
	forever := make(chan bool)

	// 5. Start a Goroutine to process messages
	go func() {
		for d := range msgs {
			// A. Decode the JSON
			var msg domain.Message
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				log.Printf("Error decoding JSON: %v", err)
				d.Nack(false, false) // Negative Ack (Reject)
				continue
			}

			// B. "Simulate" Saving to Database (The heavy work)
			log.Printf(" [x] Received from User %s: %s", msg.UserID, msg.Content)
			log.Println(" [.] Saving to DB...")
			time.Sleep(1 * time.Second) // Simulate DB latency

			// C. Manual Acknowledge (The most important part!)
			// We tell RabbitMQ: "I have finished processing. You can delete this now."
			d.Ack(false)
			log.Println(" [✓] Saved and Acknowledged")
		}
	}()

	<-forever // Block here so the main function doesn't exit
}
