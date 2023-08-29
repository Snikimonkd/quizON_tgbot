package repository

import (
	"context"
	"errors"
	"fmt"
	"quizon_bot/internal/app/usecase"
	"quizon_bot/internal/generated/postgres/public/model"
	"quizon_bot/internal/generated/postgres/public/table"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/jackc/pgx/v5"
)

type repository struct {
	db *pgx.Conn
}

func NewRepository(db *pgx.Conn) repository {
	return repository{
		db: db,
	}
}

func (r repository) Games(ctx context.Context) ([]model.Games, error) {
	stmt := table.Games.SELECT(
		table.Games.ID,
		table.Games.CreatedAt,
		table.Games.UdpatedAt,
		table.Games.Location,
		table.Games.Date,
		table.Games.Description,
	).WHERE(
		table.Games.Date.GT(postgres.TimestampzT(time.Now())),
	)

	query, args := stmt.Sql()
	rows, err := r.db.Query(ctx, query, args...)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("can't query from games: %w", err)
	}

	var res []model.Games
	for rows.Next() {
		var row model.Games
		err = rows.Scan(
			&row.ID,
			&row.CreatedAt,
			&row.UdpatedAt,
			&row.Location,
			&row.Date,
			&row.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("can't scan result of query from games: %w", err)
		}

		res = append(res, row)
	}

	return res, nil
}

func (r repository) Create(ctx context.Context, req model.Games) error {
	stmt := table.Games.INSERT(
		table.Games.AllColumns.Except(table.Games.ID),
	).MODEL(req)

	query, args := stmt.Sql()
	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("can't insert into games: %w", err)
	}

	return nil
}

func (r repository) UpsertAdmin(ctx context.Context, req model.Admins) error {
	stmt := table.Admins.INSERT(
		table.Admins.AllColumns,
	).MODEL(
		req,
	).ON_CONFLICT(
		table.Admins.ID,
	).DO_UPDATE(
		postgres.SET(
			table.Admins.DateUntil.SET(table.Admins.EXCLUDED.DateUntil),
		),
	)

	query, args := stmt.Sql()
	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("can't upsert admin: %w", err)
	}

	return nil
}

func (r repository) CheckAuth(ctx context.Context, userID int64) (model.Admins, error) {
	stmt := table.Admins.SELECT(
		table.Admins.ID,
		table.Admins.DateUntil,
	).WHERE(
		table.Admins.ID.EQ(postgres.Int64(userID)),
	)

	query, args := stmt.Sql()
	var res model.Admins
	err := r.db.QueryRow(ctx, query, args...).Scan(&res.ID, &res.DateUntil)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Admins{}, usecase.ErrNotFound
	}
	if err != nil {
		return model.Admins{}, fmt.Errorf("can't select from admins: %w", err)
	}

	return res, nil
}

func (r repository) Register(ctx context.Context, req model.Registrations) error {
	stmt := table.Registrations.INSERT(
		table.Registrations.AllColumns,
	).MODEL(
		req,
	).ON_CONFLICT(table.Registrations.GameID, table.Registrations.TeamID).
		DO_NOTHING()

	query, args := stmt.Sql()
	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("can't insert into registrations: %w", err)
	}

	return nil
}

func (r repository) List(ctx context.Context, gameID int) ([]model.Registrations, error) {
	stmt := table.Registrations.SELECT(
		table.Registrations.TeamName,
		table.Registrations.TeamID,
	).WHERE(
		table.Registrations.GameID.EQ(postgres.Int(int64(gameID))),
	).ORDER_BY(
		table.Registrations.CreatedAt.ASC(),
	)

	query, args := stmt.Sql()
	rows, err := r.db.Query(ctx, query, args...)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("can't select from registrations: %w", err)
	}

	var res []model.Registrations
	for rows.Next() {
		var row model.Registrations
		err = rows.Scan(
			&row.TeamName,
			&row.TeamID,
		)
		if err != nil {
			return nil, fmt.Errorf("can't scan result of select from registrations: %w", err)
		}
		res = append(res, row)
	}

	return res, nil
}
