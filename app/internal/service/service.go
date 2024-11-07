package service

import (
	"context"

	"github.com/kms-qwe/chat-server/internal/model"
)

// ChatService interface for service layer
type ChatService interface {
	CreateChat(ctx context.Context, usernames []string) (int64, error)
	DeleteChat(ctx context.Context, chatID int64) error
	SendMessage(cxt context.Context, message *model.Message) error
}
