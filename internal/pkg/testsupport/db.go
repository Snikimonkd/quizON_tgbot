package testsupport

import (
	"context"
	"quizon_bot/internal/config"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	db *pgxpool.Pool
	s  sync.Once
)

// ConnectToTestPostgres - подключиться к базе в тестах
func ConnectToTestPostgres() *pgxpool.Pool {
	s.Do(func() {
		db = config.ConnectToPostgres(context.Background())
	})

	return db
}
