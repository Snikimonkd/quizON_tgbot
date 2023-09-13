package testsupport

import (
	"context"
	"quizon_bot/internal/config"
	"sync"

	"github.com/jackc/pgx/v5"
)

var (
	db *pgx.Conn
	s  sync.Once
)

// ConnectToTestPostgres - подключиться к базе в тестах
func ConnectToTestPostgres() *pgx.Conn {
	s.Do(func() {
		db = config.ConnectToPostgres(context.Background())
	})

	return db
}
