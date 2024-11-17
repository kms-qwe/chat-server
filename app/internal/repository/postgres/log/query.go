package log

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

const (
	tableName = "log"

	idColumn      = "id"
	messageColumn = "message"
	logTimeColumn = "log_time"
)

func queryLog(_ context.Context, log string) (string, []interface{}, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(messageColumn).
		Values(log)

	return builder.ToSql()
}
