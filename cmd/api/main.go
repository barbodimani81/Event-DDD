package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/barbodimani81/Event-DDD.git/internal/domain"
	"github.com/barbodimani81/Event-DDD.git/internal/infra/rabbit"
	r "github.com/barbodimani81/Event-DDD.git/internal/infra/redis"
	"github.com/barbodimani81/Event-DDD.git/internal/service"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	rateLimiter := r.NewRedisRateLimiter(rdb)
	eventPublisher, err := rabbit.NewRabbitPublisher(conn)
	if err != nil {
		log.Fatalf("cannot connect to rabbit %v", err)
	}

	msgService := service.NewMessageService(rateLimiter, eventPublisher)

	http.HandleFunc("/v1/send-message", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req domain.Message
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Call the Service
		// We create a context with a timeout so the request doesn't hang forever
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		err := msgService.SendMessage(ctx, req)
		if err != nil {
			if err.Error() == "rate limit exceeded" {
				http.Error(w, "Too Many Requests", 429) // The specific status code
				return
			}
			log.Printf("Internal Error: %v", err)
			http.Error(w, "Internal Server Error", 500)
			return
		}

		// Success!
		w.WriteHeader(http.StatusAccepted) // 202 Accepted
		w.Write([]byte(`{"status": "queued"}`))
	})

	// 6. Start Server
	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
