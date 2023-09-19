package usecase

import (
	"context"
	"errors"
	"quizon_bot/internal/generated/postgres/public/model"
	"time"
)

type CheckAuthRepository interface {
	CheckAuth(ctx context.Context, userID int64) (model.Admins, error)
}

func (u usecase) CheckAuth(ctx context.Context, userID int64) (bool, error) {
	until, err := u.checkAuthRepository.CheckAuth(ctx, userID)
	if errors.Is(err, ErrNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	now := time.Now()
	if until.DateUntil.After(now) {
		return true, nil
	}

	return false, nil
}
