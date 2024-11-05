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

func (p *pgStorage) queryCreateChat(ctx context.Context, usernames []string) (*sq.InsertBuilder, int64, error) {
	if len(usernames) == 0 {
		return nil, 0, errors.New("no usernames provided for insertion")
	}

	builderChatInsert := sq.Insert("chat").
		Columns("id").
		Values(sq.Expr("DEFAULT")).
		Suffix("RETURNING id")

	query, args, err := builderChatInsert.ToSql()
	if err != nil {
		return nil, 0, err
	}
	var chatID int64
	err = p.pool.QueryRow(ctx, query, args...).Scan(&chatID)
	if err != nil {
		return nil, 0, err
	}

	builderInsert := sq.Insert("chat_to_user").
		PlaceholderFormat(sq.Dollar).
		Columns("chat_id", "user_name")

	for i := range usernames {
		builderInsert = builderInsert.Values(chatID, usernames[i])
	}
	return &builderInsert, chatID, nil
}
func (p *pgStorage) CreateChat(ctx context.Context, usernames []string) (int64, error) {
	builderInsert, chatID, err := p.queryCreateChat(ctx, usernames)
	if err != nil {
		return 0, err
	}
	query, args, err := builderInsert.ToSql()
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
	builderDelete := sq.Delete("chat").
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
	builderInsert := sq.Insert("message").
		PlaceholderFormat(sq.Dollar).
		Columns("user_name", "message_text", "chat_id", "message_time_send").
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
