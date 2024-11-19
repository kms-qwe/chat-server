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

// NewChatService creates new a UserService with provided  UserRepository LogRepository TxManager
func NewChatService(
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

// CreateChat creates chat with given participants and return its id
func (s *serv) CreateChat(ctx context.Context, usernames []string) (int64, error) {

	var chatID int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error

		chatID, errTx = s.chatRepository.CreateChat(ctx)
		if errTx != nil {
			return errTx
		}

		errTx = s.chatRepository.CreateParticipants(ctx, chatID, usernames)
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

	return chatID, nil
}

// DeleteChat deletes chat with id equal chatID
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

// SendMessage creates message with given params
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
