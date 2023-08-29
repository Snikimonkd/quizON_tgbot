package usecase

import (
	"quizon_bot/internal/generated/postgres/public/model"

	"context"
)

type GamesRepository interface {
	Games(ctx context.Context) ([]model.Games, error)
}

func (u usecase) Games(ctx context.Context) ([]model.Games, error) {
	return u.gamesRepository.Games(ctx)
}
