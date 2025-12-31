# Go Notificator

A high-performance microservice prototype demonstrating **Domain-Driven Design (DDD)**, **Rate Limiting**, and **Event-Driven Architecture**.

## 🚀 Tech Stack

* **Language:** Golang
* **Caching:** Redis (for Rate Limiting)
* **Messaging:** RabbitMQ (for Asynchronous Processing)

## 🛠️ Setup & Run

### 1. Start Infrastructure

Run Redis and RabbitMQ (using Docker is easiest):

```bash
docker run -d -p 6379:6379 redis
docker run -d -p 5672:5672 -p 15672:15672 rabbitmq:management

```

### 2. Run the Application

```bash
go mod tidy
go run cmd/api/main.go

```

## 📡 Usage

**Send a Message (POST)**

```bash
curl -X POST http://localhost:8080/v1/send-message \
     -H "Content-Type: application/json" \
     -d '{"user_id": "123", "content": "Hello Navatel!"}'

```

* **Success:** `202 Accepted`
* **Rate Limit Exceeded:** `429 Too Many Requests` (after 5 requests)

## 📂 Structure

```text
├── internal
│   ├── domain          # Pure Interfaces & Entities
│   ├── service         # Business Logic
│   └── infrastructure  # Redis & RabbitMQ implementations
└── cmd/api             # Main Entry Point

```
