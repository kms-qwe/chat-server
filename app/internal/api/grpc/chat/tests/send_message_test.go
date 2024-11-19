package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/kms-qwe/chat-server/internal/api/grpc/chat"
	"github.com/kms-qwe/chat-server/internal/model"
	"github.com/kms-qwe/chat-server/internal/service"
	serviceMocks "github.com/kms-qwe/chat-server/internal/service/mocks"
	desc "github.com/kms-qwe/chat-server/pkg/chat_v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestSendMessage(t *testing.T) {
	t.Parallel()
	type chatServiceMockFunc func(mc *minimock.Controller) service.ChatService

	type args struct {
		ctx context.Context
		req *desc.SendMessageRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id       = gofakeit.Int64()
		from     = gofakeit.Name()
		text     = gofakeit.Sentence(10)
		sendtime = gofakeit.Date()

		reqCorrect = &desc.SendMessageRequest{
			Message: &desc.Message{
				ChatId:   id,
				From:     from,
				Text:     text,
				SendTime: timestamppb.New(sendtime),
			},
		}

		serviceMessageCorrect = &model.Message{
			From:     from,
			Text:     text,
			ChatID:   id,
			SendTime: sendtime,
		}

		reqEmpty = &desc.SendMessageRequest{
			Message: &desc.Message{},
		}

		serviceMessageEmpty = &model.Message{
			SendTime: time.Unix(0, 0).UTC(),
		}

		res = &emptypb.Empty{}

		chatServiceErr = fmt.Errorf("chat service error")
	)

	tests := []struct {
		name            string
		args            args
		want            *emptypb.Empty
		err             error
		chatServiceMock chatServiceMockFunc
	}{
		{
			name: "t1: succes case",
			args: args{
				ctx: ctx,
				req: reqCorrect,
			},
			want: res,
			err:  nil,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.SendMessageMock.Expect(ctx, serviceMessageCorrect).Return(nil)
				return mock
			},
		},
		{
			name: "t2: service error",
			args: args{
				ctx: ctx,
				req: reqEmpty,
			},
			want: res,
			err:  chatServiceErr,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.SendMessageMock.Expect(ctx, serviceMessageEmpty).Return(chatServiceErr)
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

			res, err := api.SendMessage(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)

		})
	}

}
