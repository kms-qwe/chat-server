package postgres

import (
	"context"
	"errors"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kms-qwe/chat-server/internal/storage"
)

type pgStorage struct {
	pool *pgxpool.Pool
}

// NewPgStorage initializes a new PostgreSQL storage instance using the provided DSN.
func NewPgStorage(ctx context.Context, DSN string) (storage.Storage, error) {
	pool, err := pgxpool.New(ctx, DSN)
	if err != nil {
		return nil, err
	}
	return &pgStorage{
		pool: pool,
	}, nil
}

func (p *pgStorage) CreateChat(ctx context.Context, usernames []string) (int64, error) {
	if len(usernames) == 0 {
		return 0, errors.New("no usernames provided for insertion")
	}

	builderChatInsert := sq.Insert("chatV1.chat").
		Columns("id").
		Values(sq.Expr("DEFAULT")).
		Suffix("RETURNING id")

	query, args, err := builderChatInsert.ToSql()
	if err != nil {
		return 0, err
	}
	var chatID int64
	err = p.pool.QueryRow(ctx, query, args...).Scan(&chatID)
	if err != nil {
		return 0, err
	}

	builderInsert := sq.Insert("chatV1.chat_to_user").
		PlaceholderFormat(sq.Dollar).
		Columns("chat_id", "user_name")

	for i := range usernames {
		builderInsert = builderInsert.Values(chatID, usernames[i])
	}

	query, args, err = builderInsert.ToSql()
	if err != nil {
		return 0, err
	}
	_, err = p.pool.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	log.Printf("chat create with id: %d\n", chatID)

	return chatID, nil
}

func (p *pgStorage) DeleteChat(ctx context.Context, chatID int64) error {
	builderDelete := sq.Delete("chatV1.chat").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": chatID})
	query, args, err := builderDelete.ToSql()
	if err != nil {
		return err
	}
	res, err := p.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	log.Printf("delete chat with id: %d, rows affected: %d\n", chatID, res.RowsAffected())
	return nil
}

func (p *pgStorage) SendMessage(ctx context.Context, from, text string, chatID int64, timestamp time.Time) error {
	builderInsert := sq.Insert("chatV1.message").
		PlaceholderFormat(sq.Dollar).
		Columns("user_name", "message_text", "chat_id", "time_stamp").
		Values(from, text, chatID, timestamp).
		Suffix("RETURNING id")
	query, args, err := builderInsert.ToSql()
	if err != nil {
		return err
	}
	var messageID int64
	err = p.pool.QueryRow(ctx, query, args...).Scan(&messageID)
	if err != nil {
		return err
	}
	log.Printf("message create with id: %d\n", messageID)
	return nil
}
