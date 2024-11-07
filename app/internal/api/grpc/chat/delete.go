package chat

import (
	"context"
	"log"

	desc "github.com/kms-qwe/chat-server/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Delete handles the request for creating a new chat.
func (i *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	err := i.chatService.DeleteChat(ctx, req.Id)
	if err != nil {
		return &emptypb.Empty{}, err
	}

	log.Printf("delete chat with id: %d", req.Id)

	return &emptypb.Empty{}, nil
}
