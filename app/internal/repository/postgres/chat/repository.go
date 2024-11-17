package chat

import (
	"context"
	"fmt"
	"log"

	"github.com/kms-qwe/chat-server/internal/model"
	"github.com/kms-qwe/chat-server/internal/repository"
	pgClient "github.com/kms-qwe/platform_common/pkg/client/postgres"
)

type repo struct {
	db pgClient.Client
}

// NewChatRepository initializes a new PostgreSQL storage instance using the provided DSN.
func NewChatRepository(pgClient pgClient.Client) repository.ChatRepository {

	return &repo{
		db: pgClient,
	}
}

func (r *repo) CreateChat(ctx context.Context) (int64, error) {

	query, args, err := queryCreateChat(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to build query: %w", err)
	}

	q := pgClient.Query{
		Name:     "chat_repository.CraeteChat",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().ScanOneContext(ctx, &id, q, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to scan new chat id: %w", err)
	}

	return id, nil
}

func (r *repo) CreateParticipants(ctx context.Context, chatID int64, usernames []string) error {

	query, args, err := queryCreateParticipants(ctx, chatID, usernames)
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	q := pgClient.Query{
		Name:     "chat_repository.CreateParticipants",
		QueryRaw: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("failed to insert chat users: %w", err)
	}

	log.Printf("created chat with id: %d, with number of users: %d\n", chatID, res.RowsAffected())

	return nil
}

func (r *repo) DeleteChat(ctx context.Context, chatID int64) error {

	query, args, err := queryDeleteChat(ctx, chatID)
	if err != nil {
		return err
	}

	q := pgClient.Query{
		Name:     "chat_repository.Delete",
		QueryRaw: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("failed to delete chat: %w", err)
	}

	log.Printf("delete chat with id: %d, rows affected: %d\n", chatID, res.RowsAffected())
	return nil
}

func (r *repo) SendMessage(ctx context.Context, message *model.Message) error {

	query, args, err := queryCreateMessage(ctx, message)
	if err != nil {
		return err
	}

	q := pgClient.Query{
		Name:     "chat_repository.SendMassage",
		QueryRaw: query,
	}

	var messageID int64
	err = r.db.DB().ScanOneContext(ctx, &messageID, q, args...)
	if err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}

	log.Printf("message create with id: %d\n", messageID)

	return nil
}
