package usecase

import (
	"context"
	"quizon_bot/internal/generated/postgres/public/model"
)

type RegisterRepository interface {
	Register(ctx context.Context, req model.RegistrationsDraft) error
}

func (u usecase) Register(ctx context.Context, userID int64) error {
	now := u.clock.Now()
	req := model.RegistrationsDraft{
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := u.registerRepository.Register(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
