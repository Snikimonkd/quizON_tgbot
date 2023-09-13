package usecase

import (
	"quizon_bot/internal/generated/postgres/public/model"
	"quizon_bot/internal/pkg/timesupport"
	"time"

	"context"
)

type GamesRepository interface {
	Games(ctx context.Context, from time.Time) ([]model.Games, error)
}

func (u usecase) Games(ctx context.Context) ([]model.Games, error) {
	return u.gamesRepository.Games(ctx, time.Now().In(timesupport.LocMsk).Add(-time.Hour))
}
