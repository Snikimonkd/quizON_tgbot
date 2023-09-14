package testsupport

import (
	"context"
	"quizon_bot/internal/generated/postgres/public/model"
	"quizon_bot/internal/generated/postgres/public/table"
	"testing"

	"github.com/go-jet/jet/v2/postgres"
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

func TruncateGames(t *testing.T, db *pgx.Conn) {
	query := "TRUNCATE games CASCADE"
	_, err := db.Exec(context.Background(), query)
	if err != nil {
		t.Errorf("can't truncate games: %v", err.Error())
	}
}

func SelectGames(t *testing.T, db *pgx.Conn, id int64) model.Games {
	stmt := table.Games.SELECT(
		table.Games.AllColumns,
	).WHERE(
		table.Games.ID.EQ(postgres.Int64(id)),
	)

	var res model.Games
	query, args := stmt.Sql()
	err := db.QueryRow(context.Background(), query, args...).Scan(
		&res.ID,
		&res.CreatedAt,
		&res.UdpatedAt,
		&res.Location,
		&res.Date,
		&res.Description,
	)
	if err != nil {
		t.Errorf("can't select games: %v", err.Error())
	}

	return res
}
