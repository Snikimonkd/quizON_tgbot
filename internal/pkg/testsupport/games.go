package testsupport

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
)

func TruncateGames(t *testing.T, db *pgx.Conn) {
	query := "TRUNCATE games"
	_, err := db.Exec(context.Background(), query)
	if err != nil {
		t.Errorf("can't truncate registrations: %v", err.Error())
	}
}
