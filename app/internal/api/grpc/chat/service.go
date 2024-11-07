package chat

import (
	"github.com/kms-qwe/chat-server/internal/service"
	desc "github.com/kms-qwe/chat-server/pkg/chat_v1"
)

// Implementation represents the gRPC handlers that implement the UserV1Server interface
// and use the UserService for business logic operations.
type Implementation struct {
	desc.UnimplementedChatV1Server
	chatService service.ChatService
}

// NewImplementation creates a new instance of GRPCHandlers with the provided ChatService.
func NewImplementation(chatService service.ChatService) *Implementation {
	return &Implementation{
		chatService: chatService,
	}
}
