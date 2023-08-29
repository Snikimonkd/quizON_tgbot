package usecase

import (
	"context"
	"quizon_bot/internal/generated/postgres/public/model"
)

type ListRepository interface {
	List(ctx context.Context, gameID int) ([]model.Registrations, error)
}

func (u usecase) List(ctx context.Context, gameID int) ([]model.Registrations, error) {
	return u.listRepository.List(ctx, gameID)
}
