package main

import (
	"context"
	"log"
	"time"

	desc "github.com/kms-qwe/microservices_course_chat-server/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:9002"
)

func main() {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("failed to close connection: %v", err)
		}
	}()
	client := desc.NewChatV1Client(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	response, err := client.Create(ctx, &desc.CreateRequest{Usernames: []string{"a", "b", "c0"}})
	if err != nil {
		log.Fatalf("failed to create user by id: %v", err)
	}
	log.Printf("create resp: %#v\n", response)
}
