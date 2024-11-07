package chat

import (
	"context"
	"log"

	"github.com/kms-qwe/chat-server/internal/converter"
	desc "github.com/kms-qwe/chat-server/pkg/chat_v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

// SendMessage handles the request for creating a new chat.
func (i *Implementation) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	err := i.chatService.SendMessage(ctx, converter.ToMessageFromDesc(req.Message))
	if err != nil {
		return &emptypb.Empty{}, err
	}

	log.Printf("save message: %v", req.Message)

	return &emptypb.Empty{}, nil
}
