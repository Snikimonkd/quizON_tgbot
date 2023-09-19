package testsupport

import (
	"context"
	"quizon_bot/internal/generated/postgres/public/model"
	"quizon_bot/internal/generated/postgres/public/table"
	"testing"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/jackc/pgx/v5"
)

func SelectRegistrations(t *testing.T, db *pgx.Conn, gameID int64, teamID int64) model.Registrations {
	stmt := table.Registrations.SELECT(
		table.Registrations.AllColumns,
	).WHERE(
		table.Registrations.GameID.EQ(postgres.Int64(gameID)).AND(
			table.Registrations.TeamID.EQ(postgres.Int64(teamID)),
		),
	)

	query, args := stmt.Sql()
	var res model.Registrations
	err := db.QueryRow(context.Background(), query, args...).Scan(
		&res.GameID,
		&res.TeamID,
		&res.TeamName,
		&res.UserID,
		&res.CreatedAt,
		&res.UdpatedAt,
	)
	if err != nil {
		t.Errorf("can't select from registrations: %v", err.Error())
	}

	return res
}

func InsertRegistration(t *testing.T, db *pgx.Conn, in model.Registrations) {
	stmt := table.Registrations.INSERT(
		table.Registrations.AllColumns,
	).MODEL(
		in,
	)

	query, args := stmt.Sql()
	_, err := db.Exec(context.Background(), query, args...)
	if err != nil {
		t.Errorf("can't insert into registrations: %v", err.Error())
	}
}

func TruncateRegistrations(t *testing.T, db *pgx.Conn) {
	query := "TRUNCATE registrations"
	_, err := db.Exec(context.Background(), query)
	if err != nil {
		t.Errorf("can't truncate registrations: %v", err.Error())
	}
}
