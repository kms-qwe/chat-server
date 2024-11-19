package app

import (
	"context"
	"log"

	"github.com/kms-qwe/chat-server/internal/api/grpc/chat"
	"github.com/kms-qwe/chat-server/internal/config"
	"github.com/kms-qwe/chat-server/internal/config/env"
	"github.com/kms-qwe/chat-server/internal/repository"
	chatpg "github.com/kms-qwe/chat-server/internal/repository/postgres/chat"
	logpg "github.com/kms-qwe/chat-server/internal/repository/postgres/log"
	"github.com/kms-qwe/chat-server/internal/service"
	chatserv "github.com/kms-qwe/chat-server/internal/service/chat"
	"github.com/kms-qwe/platform_common/pkg/client/postgres"
	pg "github.com/kms-qwe/platform_common/pkg/client/postgres/pg"
	"github.com/kms-qwe/platform_common/pkg/client/postgres/transaction"
	"github.com/kms-qwe/platform_common/pkg/closer"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	pgClient       postgres.Client
	txManager      postgres.TxManager
	chatRepository repository.ChatRepository
	logRepository  repository.LogRepository

	chatService service.ChatService

	chatGrpcHandlers *chat.GrpcHandlers
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// PGConfig provides pgconfig
func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Panicf("failed to get postgres config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

// GRPCConfig provides grpc config
func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			log.Panicf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

// PGClient provides pg client
func (s *serviceProvider) PGClient(ctx context.Context) postgres.Client {
	if s.pgClient == nil {
		pgClient, err := pg.NewPgClient(ctx, s.PGConfig().DSN())

		if err != nil {
			log.Panicf("failed to create pg client: %v", err)
		}

		err = pgClient.DB().Ping(ctx)
		if err != nil {
			log.Panicf("ping error: %s", err.Error())
		}
		s.pgClient = pgClient

		closer.Add(s.pgClient.Close)
	}

	return s.pgClient
}

// TxManager provides tx manager
func (s *serviceProvider) TxManager(ctx context.Context) postgres.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.PGClient(ctx).DB())

	}

	return s.txManager
}

// ChatRepository provides ChatRepository
func (s *serviceProvider) ChatRepository(ctx context.Context) repository.ChatRepository {
	if s.chatRepository == nil {
		s.chatRepository = chatpg.NewChatRepository(s.PGClient(ctx))
	}

	return s.chatRepository
}

// LogRepository provides  LogRepository
func (s *serviceProvider) LogRepository(ctx context.Context) repository.LogRepository {
	if s.logRepository == nil {
		s.logRepository = logpg.NewLogRepository(s.PGClient(ctx))
	}

	return s.logRepository
}

// ChatService provides ChatService
func (s *serviceProvider) ChatService(ctx context.Context) service.ChatService {
	if s.chatService == nil {
		s.chatService = chatserv.NewChatService(s.ChatRepository(ctx), s.LogRepository(ctx), s.TxManager(ctx))
	}

	return s.chatService
}

// ChatGrpcHandlers provides ChatGrpcHandlers
func (s *serviceProvider) ChatGrpcHandlers(ctx context.Context) *chat.GrpcHandlers {
	if s.chatGrpcHandlers == nil {
		s.chatGrpcHandlers = chat.NewGrpcHandlers(s.ChatService(ctx))
	}

	return s.chatGrpcHandlers
}
