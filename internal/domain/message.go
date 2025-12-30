package domain

import "context"

type Message struct {
	UserID  string `json:"user_id"`
	Content string `json:"content"`
}

type RateLimiter interface {
	IsAllowed(ctx context.Context, userID string, limit int) (bool, error)
}

type EventPublisher interface {
	Publish(ctx context.Context, msg Message) error
}
