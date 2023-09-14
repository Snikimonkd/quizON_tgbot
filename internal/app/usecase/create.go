package usecase

import (
	"context"
	"quizon_bot/internal/generated/postgres/public/model"
	"time"
)

type CreateRepository interface {
	Create(ctx context.Context, req model.Games) (int64, error)
}

func (u usecase) Create(ctx context.Context, req model.Games) error {
	now := time.Now()
	req.CreatedAt = now
	req.UdpatedAt = now
	_, err := u.createRepository.Create(ctx, req)
	return err
}
