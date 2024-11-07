package chat

import (
	"context"
	"log"

	desc "github.com/kms-qwe/chat-server/pkg/chat_v1"
)

// Create handles the request for creating a new chat.
func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := i.chatService.CreateChat(ctx, req.Usernames)
	if err != nil {
		return nil, err
	}

	log.Printf("inserted chat with id: %d", id)

	return &desc.CreateResponse{Id: id}, nil
}
