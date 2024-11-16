package chat

import (
	"context"
	"fmt"

	"github.com/kms-qwe/chat-server/internal/model"
	"github.com/kms-qwe/chat-server/internal/repository"
	"github.com/kms-qwe/chat-server/internal/service"
	pgClient "github.com/kms-qwe/platform_common/pkg/client/postgres"
)

type serv struct {
	chatRepository repository.ChatRepository
	logRepository  repository.LogRepository
	txManager      pgClient.TxManager
}

// NewUserService creates new a UserService with provided  UserRepository LogRepository TxManager
func NewUserService(
	chatRepository repository.ChatRepository,
	logRepository repository.LogRepository,
	txManager pgClient.TxManager,
) service.ChatService {
	return &serv{
		chatRepository: chatRepository,
		logRepository:  logRepository,
		txManager:      txManager,
	}
}

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

// Delete delete a new user using the provided id
func (s *serv) DeleteChat(ctx context.Context, id int64) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		errTx = s.chatRepository.DeleteChat(ctx, id)
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
