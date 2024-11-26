package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/kms-qwe/chat-server/internal/repository"
	repositoryMock "github.com/kms-qwe/chat-server/internal/repository/mocks"
	"github.com/kms-qwe/chat-server/internal/service/chat"
	pgClient "github.com/kms-qwe/platform_common/pkg/client/postgres"
	pgClientMock "github.com/kms-qwe/platform_common/pkg/client/postgres/mocks"
	"github.com/stretchr/testify/require"
)

func TestDeleteChat(t *testing.T) {
	t.Parallel()
	type chatRepositoryMockFun func(mc *minimock.Controller) repository.ChatRepository
	type logRepositoryMockFunc func(mc *minimock.Controller) repository.LogRepository
	type txManagerMockFunc func(mc *minimock.Controller) pgClient.TxManager

	type args struct {
		ctx context.Context
		req int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id = gofakeit.Int64()

		chatRepoDeleteChatErr = fmt.Errorf("chat repo delete chat error")
		logRepoErr            = fmt.Errorf("log repo error")
		txManagerErr          = fmt.Errorf("tx manager error")

		reqCorrect = id
		reqEmpty   = int64(0)

		logCorrect = fmt.Sprintf("chat deleted: %d", id)
		logEmpty   = fmt.Sprintf("chat deleted: %d", 0)
	)

	tests := []struct {
		name               string
		args               args
		err                error
		chatRepositoryMock chatRepositoryMockFun
		logRepositoryMock  logRepositoryMockFunc
		txManagerMock      txManagerMockFunc
	}{
		{
			name: "t1: succes case",
			args: args{
				ctx: ctx,
				req: reqCorrect,
			},
			err: nil,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMock.NewChatRepositoryMock(mc)
				mock.DeleteChatMock.Expect(ctx, id).Return(nil)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repositoryMock.NewLogRepositoryMock(mc)
				mock.LogMock.Expect(ctx, logCorrect).Return(nil)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) pgClient.TxManager {
				mock := pgClientMock.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f pgClient.Handler) error {
					return f(ctx)
				})
				return mock
			},
		},

		{
			name: "t2: empty case",
			args: args{
				ctx: ctx,
				req: reqEmpty,
			},
			err: nil,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMock.NewChatRepositoryMock(mc)
				mock.DeleteChatMock.Expect(ctx, 0).Return(nil)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repositoryMock.NewLogRepositoryMock(mc)
				mock.LogMock.Expect(ctx, logEmpty).Return(nil)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) pgClient.TxManager {
				mock := pgClientMock.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f pgClient.Handler) error {
					return f(ctx)
				})
				return mock
			},
		},

		{
			name: "t3: chat repo error case",
			args: args{
				ctx: ctx,
				req: reqCorrect,
			},
			err: chatRepoDeleteChatErr,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMock.NewChatRepositoryMock(mc)
				mock.DeleteChatMock.Expect(ctx, id).Return(chatRepoDeleteChatErr)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repositoryMock.NewLogRepositoryMock(mc)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) pgClient.TxManager {
				mock := pgClientMock.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f pgClient.Handler) error {
					return f(ctx)
				})
				return mock
			},
		},

		{
			name: "t4: log repo error case",
			args: args{
				ctx: ctx,
				req: reqCorrect,
			},
			err: logRepoErr,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMock.NewChatRepositoryMock(mc)
				mock.DeleteChatMock.Expect(ctx, id).Return(nil)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repositoryMock.NewLogRepositoryMock(mc)
				mock.LogMock.Expect(ctx, logCorrect).Return(logRepoErr)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) pgClient.TxManager {
				mock := pgClientMock.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f pgClient.Handler) error {
					return f(ctx)
				})
				return mock
			},
		},

		{
			name: "t5: tx manager error case",
			args: args{
				ctx: ctx,
				req: reqCorrect,
			},
			err: txManagerErr,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMock.NewChatRepositoryMock(mc)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repositoryMock.NewLogRepositoryMock(mc)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) pgClient.TxManager {
				mock := pgClientMock.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f pgClient.Handler) error {
					return txManagerErr
				})
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			chatRepositoryMock := tt.chatRepositoryMock(mc)
			logRepositoryMock := tt.logRepositoryMock(mc)
			txManagerMock := tt.txManagerMock(mc)
			chatService := chat.NewChatService(chatRepositoryMock, logRepositoryMock, txManagerMock)

			err := chatService.DeleteChat(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
		})
	}

}
