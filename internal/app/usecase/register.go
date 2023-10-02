package usecase

import (
	"context"

	httpModel "quizon_bot/internal/app/delivery/http/model"
	"quizon_bot/internal/generated/postgres/public/model"
	"quizon_bot/internal/utils"

	"github.com/jackc/pgx/v5"
)

type RegisterRepository interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Commit(ctx context.Context, tx pgx.Tx) error
	Register(ctx context.Context, tx pgx.Tx, in model.Registrations) error
	GenerateTeamID(ctx context.Context, tx pgx.Tx) (int64, error)
}

func (u usecase) Register(ctx context.Context, req httpModel.Register) error {
	tx, err := u.registerRepository.Begin(ctx)
	if err != nil {
		return err
	}
	defer utils.RollBackUnlessCommitted(ctx, tx)

	if req.TeamID == nil {
		teamID, err := u.registerStatesRepository.GenerateTeamID(ctx, tx)
		if err != nil {
			return err
		}
		req.TeamID = &teamID
	}

	now := u.clock.Now()
	domainModel := model.Registrations{
		UserID:      req.UserID,
		TgContact:   req.TgContact,
		TeamID:      *req.TeamID,
		TeamName:    req.TeamName,
		CaptainName: req.CaptainName,
		Phone:       req.Pnohe,
		GroupName:   req.GroupName,
		Amount:      req.Amount,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err = u.registerRepository.Register(ctx, tx, domainModel)
	if err != nil {
		return err
	}

	err = u.registerRepository.Commit(ctx, tx)
	if err != nil {
		return err
	}

	return nil
}
