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

func TestCreateChat(t *testing.T) {
	t.Parallel()
	type chatRepositoryMockFun func(mc *minimock.Controller) repository.ChatRepository
	type logRepositoryMockFunc func(mc *minimock.Controller) repository.LogRepository
	type txManagerMockFunc func(mc *minimock.Controller) pgClient.TxManager

	type args struct {
		ctx context.Context
		req []string
	}

	usernames := []string{}
	for range gofakeit.Number(1, 10) {
		usernames = append(usernames, gofakeit.Name())
	}
	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id = gofakeit.Int64()

		chatRepoCreateChatErr         = fmt.Errorf("chat repo create chat error")
		chatRepoCreateParticipantsErr = fmt.Errorf("chat repo create participants error")
		logRepoErr                    = fmt.Errorf("log repo error")
		txManagerErr                  = fmt.Errorf("tx manager error")

		reqCorrect = usernames
		reqEmpty   []string

		logCorrect = fmt.Sprintf("create chat: %#v", usernames)
	)

	tests := []struct {
		name               string
		args               args
		want               int64
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
			want: id,
			err:  nil,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMock.NewChatRepositoryMock(mc)
				mock.CreateChatMock.Expect(ctx).Return(id, nil)
				mock.CreateParticipantsMock.Expect(ctx, id, usernames).Return(nil)
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
			want: 0,
			err:  chatRepoCreateParticipantsErr,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMock.NewChatRepositoryMock(mc)
				mock.CreateChatMock.Expect(ctx).Return(id, nil)
				mock.CreateParticipantsMock.Expect(ctx, id, nil).Return(chatRepoCreateParticipantsErr)
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
			name: "t3: user repo error case",
			args: args{
				ctx: ctx,
				req: reqCorrect,
			},
			want: 0,
			err:  chatRepoCreateChatErr,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMock.NewChatRepositoryMock(mc)
				mock.CreateChatMock.Expect(ctx).Return(0, chatRepoCreateChatErr)
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
			want: 0,
			err:  logRepoErr,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMock.NewChatRepositoryMock(mc)
				mock.CreateChatMock.Expect(ctx).Return(id, nil)
				mock.CreateParticipantsMock.Expect(ctx, id, usernames).Return(nil)
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
			want: 0,
			err:  txManagerErr,
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

			res, err := chatService.CreateChat(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}

}
