package service

import (
	"context"
	"fmt"

	"github.com/barbodimani81/Event-DDD.git/internal/domain"
)

type MessageService struct {
	Limiter   domain.RateLimiter
	Publisher domain.EventPublisher
}

func NewMessageService(limiter domain.RateLimiter, publisher domain.EventPublisher) *MessageService {
	return &MessageService{
		Limiter:   limiter,
		Publisher: publisher,
	}
}

func (s *MessageService) SendMessage(ctx context.Context, msg domain.Message) error {
	allowed, err := s.Limiter.IsAllowed(ctx, msg.UserID, 5)
	if err != nil {
		return fmt.Errorf("cannot send message: %w", err)
	}

	if !allowed {
		return fmt.Errorf("rate limit exceeded!: %w", err)
	}

	if err := s.Publisher.Publish(ctx, msg); err != nil {
		return fmt.Errorf("cannot publish: %w", err)
	}

	return nil
}
