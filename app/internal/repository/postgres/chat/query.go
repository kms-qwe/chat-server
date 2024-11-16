package chat

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/kms-qwe/chat-server/internal/model"
	"github.com/kms-qwe/chat-server/internal/repository/postgres/chat/converter"
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

func queryCreateChat(context.Context) (string, []interface{}, error) {
	builder := sq.Insert(chatTableName).
		Columns(chatIDColumn).
		Values(sq.Expr("DEFAULT")).
		Suffix(fmt.Sprintf("RETURNING %s", chatIDColumn))

	return builder.ToSql()
}

func queryCreateParticipants(_ context.Context, chatID int64, usernames []string) (string, []interface{}, error) {
	builder := sq.Insert(chatParticipantsTableName).
		PlaceholderFormat(sq.Dollar).
		Columns(chatIDFkColumn, userNameColumn)

	for _, username := range usernames {
		builder = builder.Values(chatID, username)
	}

	return builder.ToSql()
}

func queryDeleteChat(_ context.Context, chatID int64) (string, []interface{}, error) {
	builder := sq.Delete(chatTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{chatIDColumn: chatID})

	return builder.ToSql()
}

func queryCreateMessage(_ context.Context, message *model.Message) (string, []interface{}, error) {
	repoMessage := converter.ToRepoFromMessage(message)

	builder := sq.Insert(messageTableName).
		PlaceholderFormat(sq.Dollar).
		Columns(userNameColumn, messageTextColumn, chatIDFkColumn, messageTimeSendColumn).
		Values(repoMessage.From, repoMessage.Text, repoMessage.ChatID, repoMessage.SendTime).
		Suffix(fmt.Sprintf("RETURNING %s", chatIDColumn))

	return builder.ToSql()
}
