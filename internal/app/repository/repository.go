package repository

import (
	"context"
	"errors"
	"fmt"
	"quizon_bot/internal/app/usecase"
	"quizon_bot/internal/generated/postgres/public/model"
	"quizon_bot/internal/generated/postgres/public/table"
	"quizon_bot/internal/logger"

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
	).ORDER_BY(
		table.Registrations.CreatedAt.DESC(),
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
	).WHERE(
		table.UserState.UserID.EQ(postgres.Int64(userID)),
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
	stmt := table.RegistrationsDraft.INSERT(
		table.RegistrationsDraft.AllColumns,
	).MODEL(
		in,
	).ON_CONFLICT(
		table.RegistrationsDraft.UserID,
	).DO_UPDATE(
		postgres.SET(
			table.RegistrationsDraft.TeamID.SET(table.RegistrationsDraft.EXCLUDED.TeamID),
			table.RegistrationsDraft.TeamName.SET(table.RegistrationsDraft.EXCLUDED.TeamName),
			table.RegistrationsDraft.CaptainName.SET(table.RegistrationsDraft.EXCLUDED.CaptainName),
			table.RegistrationsDraft.GroupName.SET(table.RegistrationsDraft.EXCLUDED.GroupName),
			table.RegistrationsDraft.Phone.SET(table.RegistrationsDraft.EXCLUDED.Phone),
			table.RegistrationsDraft.Amount.SET(table.RegistrationsDraft.EXCLUDED.Amount),
			table.RegistrationsDraft.UpdatedAt.SET(table.RegistrationsDraft.EXCLUDED.UpdatedAt),
		),
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

func (r repository) Start(ctx context.Context, userID int64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("can' begin tx: %w", err)
	}
	defer RollBackUnlessCommitted(ctx, tx)

	stmt := table.UserState.DELETE().WHERE(table.UserState.UserID.EQ(postgres.Int64(userID)))
	query, args := stmt.Sql()
	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("can't delete from user_states: %w", err)
	}

	stmt = table.RegistrationsDraft.DELETE().WHERE(table.RegistrationsDraft.UserID.EQ(postgres.Int64(userID)))
	query, args = stmt.Sql()
	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("can't delete from registrations_draft: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("can't commit tx: %w", err)
	}

	return nil
}
