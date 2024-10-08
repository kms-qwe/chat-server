package main

import (
	"context"
	"fmt"
	"log"
	"net"

	desc "github.com/kms-qwe/microservices_course_chat-server/app/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	grpcPort = "9002"
)

type server struct {
	desc.UnimplementedChatV1Server
}

func (s *server) Create(_ context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Printf("get create request: %#v:\n", req)
	return &desc.CreateResponse{Id: 0}, nil
}

func (s *server) Delete(_ context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("get delete request: %#v\n", req)
	return &emptypb.Empty{}, nil
}

func (s *server) SendMassage(_ context.Context, req *desc.SendMassageRequest) (*emptypb.Empty, error) {
	log.Printf("get send message request: %#v\n", req)
	return &emptypb.Empty{}, nil
}
func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen %s: %v", grpcPort, err)
	}

	serv := grpc.NewServer()
	reflection.Register(serv)
	desc.RegisterChatV1Server(serv, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = serv.Serve(lis); err != nil {
		log.Fatalf("falied to serve %v", err)
	}
}
