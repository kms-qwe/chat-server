package converter

import (
	"github.com/kms-qwe/chat-server/internal/model"
	desc "github.com/kms-qwe/chat-server/pkg/chat_v1"
)

// ToMessageFromAPI convert api model to  service model
func ToMessageFromAPI(message *desc.Message) *model.Message {
	return &model.Message{
		From:     message.From,
		Text:     message.Text,
		ChatID:   message.ChatId,
		SendTime: message.SendTime.AsTime(),
	}
}
