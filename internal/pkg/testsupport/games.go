package testsupport

import (
	"context"
	"quizon_bot/internal/generated/postgres/public/model"
	"quizon_bot/internal/generated/postgres/public/table"
	"testing"

	"github.com/jackc/pgx/v5"
)

func InsertIntoGames(t *testing.T, db *pgx.Conn, in model.Games) {
	stmt := table.Games.INSERT(
		table.Games.AllColumns,
	).MODEL(
		in,
	)

	query, args := stmt.Sql()
	_, err := db.Exec(context.Background(), query, args...)
	if err != nil {
		t.Errorf("can't insert into games: %v", err.Error())
	}
}
