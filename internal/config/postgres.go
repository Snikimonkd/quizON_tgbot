package config

import (
	"context"

	"quizon_bot/internal/logger"

	"github.com/jackc/pgx/v5"
)

// ConnectToPostgres - подключается к postgres
func ConnectToPostgres(ctx context.Context) *pgx.Conn {
	db, err := pgx.Connect(ctx, GlobalConfig.Database.DSN)
	if err != nil {
		logger.Fatalf("Can't connect to db: %v\n", err)
	}

	return db
}
