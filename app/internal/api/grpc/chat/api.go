package chat

import (
	"context"
	"log"

	"github.com/kms-qwe/chat-server/internal/converter"
	"github.com/kms-qwe/chat-server/internal/service"
	desc "github.com/kms-qwe/chat-server/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// GrpcHandlers represents the gRPC handlers that implement the UserV1Server interface
// and use the UserService for business logic operations.
type GrpcHandlers struct {
	desc.UnimplementedChatV1Server
	chatService service.ChatService
}

// NewGrpcHandlers creates a new instance of GRPCHandlers with the provided ChatService.
func NewGrpcHandlers(chatService service.ChatService) *GrpcHandlers {
	return &GrpcHandlers{
		chatService: chatService,
	}
}

// Create handles the request for creating a new chat.
func (g *GrpcHandlers) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := g.chatService.CreateChat(ctx, req.Usernames)
	if err != nil {
		return nil, err
	}

	log.Printf("inserted chat with id: %d", id)

	return &desc.CreateResponse{Id: id}, nil
}

// Delete handles the request for creating a new chat.
func (g *GrpcHandlers) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	err := g.chatService.DeleteChat(ctx, req.Id)
	if err != nil {
		return &emptypb.Empty{}, err
	}

	log.Printf("delete chat with id: %d", req.Id)

	return &emptypb.Empty{}, nil
}

// SendMessage handles the request for creating a new chat.
func (g *GrpcHandlers) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	err := g.chatService.SendMessage(ctx, converter.ToMessageFromAPI(req.Message))
	if err != nil {
		return &emptypb.Empty{}, err
	}

	log.Printf("save message: %v", req.Message)

	return &emptypb.Empty{}, nil
}
