package chat

import (
	"context"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/kms-qwe/chat-server/internal/model"
	"github.com/kms-qwe/chat-server/internal/repository"
	"github.com/kms-qwe/chat-server/internal/repository/postgres/chat/converter"
	pgClient "github.com/kms-qwe/platform_common/pkg/client/postgres"
)

const (
	chatParticipantsTableName = "chat_participants"
	chatTableName             = "chat"
	messageTableName          = "message"

	chatIDColumn          = "id"
	chatIDFkColumn        = "chat_id"
	userNameColumn        = "user_name"
	messageIDColumn       = "id"
	messageTextColumn     = "message_text"
	messageTimeSendColumn = "message_time_send"
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

func (r *repo) getNewChat(ctx context.Context) (int64, error) {
	builder := sq.Insert(chatTableName).
		Columns(chatIDColumn).
		Values(sq.Expr("DEFAULT")).
		Suffix(fmt.Sprintf("RETURNING %s", chatIDColumn))

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build query: %w", err)
	}

	q := pgClient.Query{
		Name:     "chat_repository.getNewChat",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().ScanOneContext(ctx, &id, q, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to scan new chat id: %w", err)
	}

	return id, nil
}

func (r *repo) CreateChat(ctx context.Context, usernames []string) (int64, error) {
	chatID, err := r.getNewChat(ctx)
	if err != nil {
		return 0, err
	}

	builder := sq.Insert(chatParticipantsTableName).
		PlaceholderFormat(sq.Dollar).
		Columns(chatIDFkColumn, userNameColumn)

	for _, username := range usernames {
		builder = builder.Values(chatID, username)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build query: %w", err)
	}

	q := pgClient.Query{
		Name:     "chat_repository.getNewChat",
		QueryRaw: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to insert chat users: %w", err)
	}

	log.Printf("created chat with id: %d, with number of users: %d\n", chatID, res.RowsAffected())

	return chatID, nil
}

func (r *repo) DeleteChat(ctx context.Context, chatID int64) error {
	builder := sq.Delete(chatTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{chatIDColumn: chatID})

	query, args, err := builder.ToSql()
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
	repoMessage := converter.ToRepoFromMessage(message)

	builder := sq.Insert(messageTableName).
		PlaceholderFormat(sq.Dollar).
		Columns(userNameColumn, messageTextColumn, chatIDFkColumn, messageTimeSendColumn).
		Values(repoMessage.From, repoMessage.Text, repoMessage.ChatID, repoMessage.SendTime).
		Suffix(fmt.Sprintf("RETURNING %s", chatIDColumn))

	query, args, err := builder.ToSql()
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
