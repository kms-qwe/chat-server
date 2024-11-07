package converter

import (
	"github.com/kms-qwe/chat-server/internal/model"
	desc "github.com/kms-qwe/chat-server/pkg/chat_v1"
)

// ToMessageFromDesc convert desc model to  service model
func ToMessageFromDesc(message *desc.Message) *model.Message {
	return &model.Message{
		From:      message.From,
		Text:      message.Text,
		ChatID:    message.ChatId,
		Timestamp: message.Timestamp.AsTime(),
	}
}
