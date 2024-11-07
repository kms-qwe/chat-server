package repository

import (
	"context"

	"github.com/kms-qwe/chat-server/internal/model"
)

// ChatRepository interface defines methods for user data storage operations.
type ChatRepository interface {
	CreateChat(ctx context.Context, usernames []string) (int64, error)
	DeleteChat(ctx context.Context, chatID int64) error
	SendMessage(cxt context.Context, message *model.Message) error
}

// LogRepository interface defines methods for log storage operations.
type LogRepository interface {
	Log(ctx context.Context, operation string) error
}
