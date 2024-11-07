package chat

import (
	"context"
	"fmt"
)

// Create creates a new user using the provided user model
func (s *serv) CreateChat(ctx context.Context, usernames []string) (int64, error) {

	var id int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.chatRepository.CreateChat(ctx, usernames)
		if errTx != nil {
			return errTx
		}

		errTx = s.logRepository.Log(ctx, fmt.Sprintf("create chat: %#v", usernames))
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}
