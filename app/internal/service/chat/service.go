package chat

import (
	pgClient "github.com/kms-qwe/chat-server/internal/client/postgres"
	"github.com/kms-qwe/chat-server/internal/repository"
	"github.com/kms-qwe/chat-server/internal/service"
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
