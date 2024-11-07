package chat

import (
	"context"
	"fmt"

	"github.com/kms-qwe/chat-server/internal/model"
)

// Create creates a new user using the provided user model
func (s *serv) SendMessage(ctx context.Context, message *model.Message) error {

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		errTx = s.chatRepository.SendMessage(ctx, message)
		if errTx != nil {
			return errTx
		}

		errTx = s.logRepository.Log(ctx, fmt.Sprintf("save message: %#v", *message))
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
