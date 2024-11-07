package app

import (
	"context"
	"log"

	"github.com/kms-qwe/chat-server/internal/api/grpc/chat"
	"github.com/kms-qwe/chat-server/internal/client/postgres"
	pgv1 "github.com/kms-qwe/chat-server/internal/client/postgres/pg_v1"
	"github.com/kms-qwe/chat-server/internal/client/postgres/transaction"
	"github.com/kms-qwe/chat-server/internal/closer"
	"github.com/kms-qwe/chat-server/internal/config"
	"github.com/kms-qwe/chat-server/internal/config/env"
	"github.com/kms-qwe/chat-server/internal/repository"
	chatpg "github.com/kms-qwe/chat-server/internal/repository/postgres/chat"
	logpg "github.com/kms-qwe/chat-server/internal/repository/postgres/log"
	"github.com/kms-qwe/chat-server/internal/service"
	chatserv "github.com/kms-qwe/chat-server/internal/service/chat"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	pgClient       postgres.Client
	txManager      postgres.TxManager
	chatRepository repository.ChatRepository
	logRepository  repository.LogRepository

	userService service.ChatService

	userImpl *chat.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}
func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get postgres config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) PGClient(ctx context.Context) postgres.Client {
	if s.pgClient == nil {
		pgClient, err := pgv1.NewPgClient(ctx, s.PGConfig().DSN())

		if err != nil {
			log.Fatalf("failed to create pg client: %v", err)
		}

		err = pgClient.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		s.pgClient = pgClient

		closer.Add(s.pgClient.Close)
	}

	return s.pgClient
}

func (s *serviceProvider) TxManager(ctx context.Context) postgres.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.PGClient(ctx).DB())

	}

	return s.txManager
}

func (s *serviceProvider) ChatRepository(ctx context.Context) repository.ChatRepository {
	if s.chatRepository == nil {
		s.chatRepository = chatpg.NewChatRepository(s.PGClient(ctx))
	}

	return s.chatRepository
}

func (s *serviceProvider) LogRepository(ctx context.Context) repository.LogRepository {
	if s.logRepository == nil {
		s.logRepository = logpg.NewLogRepository(s.PGClient(ctx))
	}

	return s.logRepository
}

func (s *serviceProvider) ChatService(ctx context.Context) service.ChatService {
	if s.userService == nil {
		s.userService = chatserv.NewUserService(s.ChatRepository(ctx), s.LogRepository(ctx), s.TxManager(ctx))
	}

	return s.userService
}

func (s *serviceProvider) UserImpl(ctx context.Context) *chat.Implementation {
	if s.userImpl == nil {
		s.userImpl = chat.NewImplementation(s.ChatService(ctx))
	}

	return s.userImpl
}
