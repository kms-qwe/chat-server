package log

import (
	"context"
	"fmt"

	"github.com/kms-qwe/chat-server/internal/repository"
	pgClient "github.com/kms-qwe/platform_common/pkg/client/postgres"
)

type repo struct {
	db pgClient.Client
}

// NewLogRepository create log interface
func NewLogRepository(db pgClient.Client) repository.LogRepository {
	return &repo{
		db: db,
	}
}

// Log saves log in db
func (r *repo) Log(ctx context.Context, log string) error {

	query, args, err := queryLog(ctx, log)
	if err != nil {
		return fmt.Errorf("failed to create query: %w", err)
	}

	q := pgClient.Query{
		Name:     "log_repository.Log",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("failed to create log: %w", err)
	}

	return nil
}
