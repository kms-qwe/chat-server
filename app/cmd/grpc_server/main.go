package main

import (
	"context"
	"log"

	"github.com/kms-qwe/chat-server/internal/app"
)

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %s\n", err.Error())
	}

	err = a.Run(ctx)
	if err != nil {
		log.Fatalf("failed to run app: %s\n", err.Error())
	}
}
