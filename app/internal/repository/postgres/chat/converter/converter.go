package converter

import (
	"github.com/kms-qwe/chat-server/internal/model"
	modelRepo "github.com/kms-qwe/chat-server/internal/repository/postgres/chat/model"
)

// ToMessageFromRepo convert serivce model to repo model
func ToMessageFromRepo(message *modelRepo.Message) *model.Message {
	return &model.Message{
		From:      message.From,
		Text:      message.From,
		ChatID:    message.ChatID,
		Timestamp: message.Timestamp,
	}
}

// ToRepoFromMessage convert serivce model to repo model
func ToRepoFromMessage(message *model.Message) *modelRepo.Message {
	return &modelRepo.Message{
		From:      message.From,
		Text:      message.From,
		ChatID:    message.ChatID,
		Timestamp: message.Timestamp,
	}
}
