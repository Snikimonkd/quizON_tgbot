package usecase

import (
	"context"
	"quizon_bot/internal/generated/postgres/public/model"
	"time"
)

type LoginRepository interface {
	UpsertAdmin(ctx context.Context, req model.Admins) error
}

func (u usecase) Login(ctx context.Context, userID int64, secretKey string) (bool, error) {
	if secretKey != "a3abe55d-9327-4ef4-ad0c-6d4fe94e85ec" {
		return false, nil
	}

	until := time.Now().Add(time.Hour * 24)
	req := model.Admins{
		ID:        userID,
		DateUntil: until,
	}
	err := u.loginRepository.UpsertAdmin(ctx, req)
	if err != nil {
		return false, err
	}

	return true, nil
}
