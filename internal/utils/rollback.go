package utils

import (
	"context"
	"quizon_bot/internal/logger"

	"github.com/jackc/pgx/v5"
)

// RollBackUnlessCommitted - роллбэк, если транзакция не закоммичена
func RollBackUnlessCommitted(ctx context.Context, tx pgx.Tx) {
	if tx == nil {
		return
	}

	err := tx.Rollback(ctx)
	if err == pgx.ErrTxClosed {
		return
	}

	if err != nil {
		logger.Errorf("can't rollback transaction: %w", err)
	}
}
