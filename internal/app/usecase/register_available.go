package usecase

import (
	"context"
)

type RegisterAvailableRepository interface {
	RegisterAvailable(ctx context.Context) (bool, error)
}

func (u usecase) RegisterAvailable(ctx context.Context) (bool, error) {
	ok, err := u.registerAvailableRepository.RegisterAvailable(ctx)
	if err != nil {
		return false, err
	}

	return ok, nil
}
