package storage

import (
	"context"
	"time"
)

// Storage interface defines methods for user data storage operations.
type Storage interface {
	CreateChat(ctx context.Context, usernames []string) (int64, error)
	DeleteChat(ctx context.Context, chatID int64) error
	SendMessage(cxt context.Context, from, text string, chatID int64, timestamp time.Time) error
}
