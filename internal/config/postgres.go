package config

import (
	"context"

	"quizon_bot/internal/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ConnectToPostgres - подключается к postgres
func ConnectToPostgres(ctx context.Context) *pgxpool.Pool {
	db, err := pgxpool.New(ctx, GlobalConfig.Database.DSN)
	if err != nil {
		logger.Fatalf("Can't connect to db: %v\n", err)
	}

	return db
}
