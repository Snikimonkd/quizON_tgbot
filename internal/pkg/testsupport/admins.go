package testsupport

import (
	"context"
	"quizon_bot/internal/generated/postgres/public/model"
	"quizon_bot/internal/generated/postgres/public/table"
	"testing"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/jackc/pgx/v5"
)

func SelectAdmins(t *testing.T, db *pgx.Conn, id int64) model.Admins {
	stmt := table.Admins.SELECT(
		table.Admins.AllColumns,
	).WHERE(
		table.Admins.ID.EQ(postgres.Int64(id)),
	)

	var res model.Admins
	query, args := stmt.Sql()
	err := db.QueryRow(context.Background(), query, args...).Scan(
		&res.ID,
		&res.DateUntil,
	)
	if err != nil {
		t.Errorf("can't select from admins: %v", err.Error())
	}

	return res
}

func InsertIntoAdmins(t *testing.T, db *pgx.Conn, in model.Admins) {
	stmt := table.Admins.INSERT(
		table.Admins.AllColumns,
	).MODEL(
		in,
	)

	query, args := stmt.Sql()
	_, err := db.Exec(context.Background(), query, args...)
	if err != nil {
		t.Errorf("can't insert into admins: %v", err.Error())
	}
}

func TruncateAdmins(t *testing.T, db *pgx.Conn) {
	query := "TRUNCATE admins;"
	_, err := db.Exec(context.Background(), query, )
	if err != nil {
		t.Errorf("can't insert into admins: %v", err.Error())
	} 
}
