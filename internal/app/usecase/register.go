package usecase

import (
	"context"
	"quizon_bot/internal/generated/postgres/public/model"
	"time"
)

type RegisterRepository interface {
	Register(ctx context.Context, req model.Registrations) error
}

func (u usecase) Register(ctx context.Context, req model.Registrations) error {
	req.CreatedAt = time.Now()
	err := u.registerRepository.Register(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
