package repository

import (
	"context"
	"errors"
	"fmt"
	"quizon_bot/internal/app/usecase"
	"quizon_bot/internal/generated/postgres/public/model"
	"quizon_bot/internal/generated/postgres/public/table"
	"quizon_bot/internal/logger"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"
)

type repository struct {
	db *pgx.Conn
}

func NewRepository(db *pgx.Conn) repository {
	return repository{
		db: db,
	}
}

// RollBackUnlessCommitted - роллбэк, если транзакция не закоммичена
func RollBackUnlessCommitted(ctx context.Context, tx pgx.Tx) {
	if tx == nil {
		return
	}

	err := tx.Rollback(ctx)
	if err == pgx.ErrTxClosed {
		return
	}

	if err != nil {
		logger.Errorf("can't rollback transaction: %w", err)
	}
}

func (r repository) Registrations(ctx context.Context) ([]model.Registrations, error) {
	stmt := table.Registrations.SELECT(
		table.Registrations.AllColumns,
	)

	query, args := stmt.Sql()
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("can't select registrations: %w", err)
	}
	defer rows.Close()

	var res []model.Registrations
	for rows.Next() {
		var buf model.Registrations
		err := rows.Scan(
			&buf.UserID,
			&buf.TgContact,
			&buf.TeamID,
			&buf.TeamName,
			&buf.CaptainName,
			&buf.Pnohe,
			&buf.GroupName,
			&buf.Amount,
			&buf.CreatedAt,
			&buf.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("can't scan result: %w", err)
		}

		res = append(res, buf)
	}

	return res, nil
}

func (r repository) GetState(ctx context.Context, userID int64) (string, error) {
	stmt := table.UserState.SELECT(
		table.UserState.State,
	)

	query, args := stmt.Sql()
	var state string
	err := r.db.QueryRow(ctx, query, args...).Scan(&state)
	if errors.Is(err, pgx.ErrNoRows) {
		return string(usecase.EMPTY), nil
	}
	if err != nil {
		return "", fmt.Errorf("can't get user state: %w", err)
	}

	return state, nil
}

func (r repository) UpdateState(ctx context.Context, state model.UserState) error {
	stmt := table.UserState.INSERT(
		table.UserState.AllColumns,
	).MODEL(
		state,
	).ON_CONFLICT(
		table.UserState.UserID,
	).DO_UPDATE(
		postgres.SET(
			table.UserState.State.SET(table.UserState.EXCLUDED.State),
		),
	)

	query, args := stmt.Sql()
	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("can't update user state: %w", err)
	}

	return nil
}

func (r repository) CheckTeamsAmount(ctx context.Context) (int64, error) {
	stmt := table.Registrations.SELECT(
		postgres.COUNT(postgres.STAR),
	)

	query, args := stmt.Sql()
	var res *int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&res)
	if err != nil {
		return 0, fmt.Errorf("can't select max teams amount: %w", err)
	}

	return lo.FromPtr(res), nil
}

func (r repository) Games(ctx context.Context, from time.Time) ([]model.Games, error) {
	stmt := table.Games.SELECT(
		table.Games.ID,
		table.Games.CreatedAt,
		table.Games.UdpatedAt,
		table.Games.Location,
		table.Games.Date,
		table.Games.Description,
	).WHERE(
		table.Games.Date.GT_EQ(postgres.TimestampzT(from)),
	).ORDER_BY(
		table.Games.Date.ASC(),
	)

	query, args := stmt.Sql()
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("can't query from games: %w", err)
	}
	defer rows.Close()

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

func (r repository) Create(ctx context.Context, req model.Games) (int64, error) {
	stmt := table.Games.INSERT(table.Games.AllColumns.Except(table.Games.ID)).
		MODEL(req).
		RETURNING(table.Games.ID)

	query, args := stmt.Sql()
	var res int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&res)
	if err != nil {
		return 0, fmt.Errorf("can't insert into games: %w", err)
	}

	return res, nil
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

func (r repository) DeletDraft(ctx context.Context, userID int64) error {
	stmt := table.RegistrationsDraft.DELETE().WHERE(
		table.Registrations.UserID.EQ(postgres.Int64(userID)),
	)

	query, args := stmt.Sql()
	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("can't delete from registrations_draft: %w", err)
	}

	return nil
}

func (r repository) RegisterStart(ctx context.Context, req model.RegistrationsDraft) error {
	stmt := table.RegistrationsDraft.INSERT(
		table.RegistrationsDraft.AllColumns,
	).MODEL(
		req,
	)

	query, args := stmt.Sql()
	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("can't insert into registrations_draft: %w", err)
	}

	return nil
}

// func (r repository) List(ctx context.Context, gameID int) ([]model.Registrations, error) {
// 	stmt := table.Registrations.SELECT(
// 		table.Registrations.TeamName,
// 		table.Registrations.TeamID,
// 	).WHERE(
// 		table.Registrations.GameID.EQ(postgres.Int(int64(gameID))),
// 	).ORDER_BY(
// 		table.Registrations.CreatedAt.ASC(),
// 	)
//
// 	query, args := stmt.Sql()
// 	rows, err := r.db.Query(ctx, query, args...)
// 	if err != nil {
// 		return nil, fmt.Errorf("can't select from registrations: %w", err)
// 	}
// 	defer rows.Close()
//
// 	var res []model.Registrations
// 	for rows.Next() {
// 		var row model.Registrations
// 		err = rows.Scan(
// 			&row.TeamName,
// 			&row.TeamID,
// 		)
// 		if err != nil {
// 			return nil, fmt.Errorf("can't scan result of select from registrations: %w", err)
// 		}
// 		res = append(res, row)
// 	}
//
// 	return res, nil
// }

func (r repository) GetRegistrationDraft(ctx context.Context, userID int64) (model.RegistrationsDraft, error) {
	stmt := table.RegistrationsDraft.SELECT(
		table.RegistrationsDraft.AllColumns,
	).WHERE(
		table.RegistrationsDraft.UserID.EQ(postgres.Int64(userID)),
	)

	query, args := stmt.Sql()
	var res model.RegistrationsDraft
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&res.UserID,
		&res.TgContact,
		&res.TeamID,
		&res.TeamName,
		&res.CaptainName,
		&res.GroupName,
		&res.Phone,
		&res.Amount,
		&res.CreatedAt,
		&res.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return model.RegistrationsDraft{}, usecase.ErrNotFound
	}
	if err != nil {
		return model.RegistrationsDraft{}, fmt.Errorf("can't select from registrations draft: %w", err)
	}

	return res, nil
}

func (r repository) UpdateRegistrationDraft(ctx context.Context, in model.RegistrationsDraft) error {
	stmt := table.RegistrationsDraft.UPDATE(
		table.RegistrationsDraft.TeamID,
		table.RegistrationsDraft.TeamName,
		table.RegistrationsDraft.CaptainName,
		table.RegistrationsDraft.GroupName,
		table.RegistrationsDraft.Phone,
		table.RegistrationsDraft.Amount,
		table.RegistrationsDraft.UpdatedAt,
	).SET(
		in.TeamID,
		in.TeamName,
		in.CaptainName,
		in.GroupName,
		in.Phone,
		in.Amount,
		in.UpdatedAt,
	).WHERE(
		table.RegistrationsDraft.UserID.EQ(postgres.Int64(in.UserID)),
	)

	query, args := stmt.Sql()
	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("can't update registrations_draft: %w", err)
	}

	return nil
}

func (r repository) GenerateTeamID(ctx context.Context) (int64, error) {
	query := `SELECT nextval('team_id_seq');`
	var res int64
	err := r.db.QueryRow(ctx, query).Scan(&res)
	if err != nil {
		return 0, fmt.Errorf("can't generate team_id: %w", err)
	}

	return res, nil
}

func (r repository) CheckGameID(ctx context.Context, gameID int64) (bool, error) {
	stmt := table.Games.SELECT(postgres.Int64(1)).WHERE(table.Games.ID.EQ(postgres.Int64(gameID)))
	query, args := stmt.Sql()
	var res int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&res)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("can't check game_id: %w", err)
	}

	return true, nil
}

func (r repository) CreateRegistration(ctx context.Context, in model.Registrations) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("can' begin tx: %w", err)
	}
	defer RollBackUnlessCommitted(ctx, tx)

	createRegStmt := table.Registrations.INSERT(
		table.Registrations.AllColumns,
	).MODEL(
		in,
	)

	query, args := createRegStmt.Sql()
	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("can't insert into registrations: %w", err)
	}

	deleteDraftStmt := table.RegistrationsDraft.DELETE().
		WHERE(
			table.RegistrationsDraft.UserID.EQ(postgres.Int64(in.UserID)),
		)

	query, args = deleteDraftStmt.Sql()
	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("can't delete registration draft: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("can't commit tx: %w", err)
	}

	return nil
}
