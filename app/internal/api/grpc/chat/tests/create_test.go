package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/kms-qwe/chat-server/internal/api/grpc/chat"
	"github.com/kms-qwe/chat-server/internal/service"
	serviceMocks "github.com/kms-qwe/chat-server/internal/service/mocks"
	desc "github.com/kms-qwe/chat-server/pkg/chat_v1"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	type chatServiceMockFunc func(mc *minimock.Controller) service.ChatService

	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}

	usernames := []string{}
	for range gofakeit.Number(1, 10) {
		usernames = append(usernames, gofakeit.Name())
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id = gofakeit.Int64()

		reqCorrect = &desc.CreateRequest{
			Usernames: usernames,
		}

		resCorrect = &desc.CreateResponse{
			Id: id,
		}

		reqEmpty = &desc.CreateRequest{
			Usernames: nil,
		}

		chatServiceErr = fmt.Errorf("chat service error")
	)

	tests := []struct {
		name            string
		args            args
		want            *desc.CreateResponse
		err             error
		chatServiceMock chatServiceMockFunc
	}{
		{
			name: "t1: succes case",
			args: args{
				ctx: ctx,
				req: reqCorrect,
			},
			want: resCorrect,
			err:  nil,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.CreateChatMock.Expect(ctx, usernames).Return(id, nil)
				return mock
			},
		},
		{
			name: "t2: service error",
			args: args{
				ctx: ctx,
				req: reqEmpty,
			},
			want: nil,
			err:  chatServiceErr,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.CreateChatMock.Expect(ctx, nil).Return(0, chatServiceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			chatServiceMock := tt.chatServiceMock(mc)
			api := chat.NewGrpcHandlers(chatServiceMock)

			res, err := api.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)

		})
	}

}
