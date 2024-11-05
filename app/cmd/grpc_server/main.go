package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/kms-qwe/chat-server/internal/config"
	"github.com/kms-qwe/chat-server/internal/config/env"
	"github.com/kms-qwe/chat-server/internal/storage"
	"github.com/kms-qwe/chat-server/internal/storage/postgres"
	desc "github.com/kms-qwe/chat-server/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedChatV1Server
	storage storage.Storage
}

// NewServer initializes a new server instance with the provided DSN for storage.
func NewServer(ctx context.Context, DSN string) (*server, error) {
	storage, err := postgres.NewPgStorage(ctx, DSN)
	if err != nil {
		return nil, err
	}
	return &server{
		storage: storage,
	}, nil
}
func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {

	id, err := s.storage.CreateChat(ctx, req.GetUsernames())
	if err != nil {
		log.Printf("create chat request error: %#v\n", err)
		return nil, err
	}
	return &desc.CreateResponse{Id: id}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	err := s.storage.DeleteChat(ctx, req.GetId())
	if err != nil {
		log.Printf("delete request error: %#v\n", err)
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	err := s.storage.SendMessage(ctx, req.Message.GetFrom(), req.Message.GetText(), req.Message.GetChatId(), req.Message.Timestamp.AsTime())
	if err != nil {
		log.Printf("send message request error: %#v\n", err)
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

func main() {
	flag.Parse()
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %#v\n", err)
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %#v\n", err)
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %#v\n", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to get tcp listener: %#v\n", err)
	}

	serv, err := NewServer(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to get serv: %#v\n", err)
	}
	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterChatV1Server(s, serv)

	log.Printf("server listening at %v", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %#v\n", err)
	}
}
