package usecase

import (
	"context"
	"errors"
	"quizon_bot/internal/generated/postgres/public/model"
	"time"
)

type AuthRepository interface {
	CheckAuth(ctx context.Context, userID int64) (model.Admins, error)
}

func (u usecase) CheckAuth(ctx context.Context, userID int64) (bool, error) {
	auth, err := u.authRepository.CheckAuth(ctx, userID)
	if errors.Is(err, ErrNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	if auth.DateUntil.Before(time.Now()) {
		return false, nil
	}

	return true, nil
}
