package chat

import (
	"context"
	"errors"
	"fmt"
)

// Delete delete a new user using the provided id
func (s *serv) DeleteChat(ctx context.Context, id int64) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		errTx = s.chatRepository.DeleteChat(ctx, id)
		if errTx != nil {
			return errTx
		}

		errTx = errors.New("fail on delete chat")
		if errTx != nil {
			return errTx
		}

		errTx = s.logRepository.Log(ctx, fmt.Sprintf("chat deleted: %d", id))
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
