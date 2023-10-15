package usecase

import (
	"context"

	httpModel "quizon_bot/internal/app/delivery/http/model"
	"quizon_bot/internal/generated/postgres/public/model"
)

type RegisterAvailableRepository interface {
	RegisterAvailable(ctx context.Context) (bool, error)
	SelectRegistrationRestrictions(ctx context.Context) (model.Games, error)
	RegistrationsAmount(ctx context.Context) (int64, error)
}

func (u usecase) RegisterAvailable(ctx context.Context) (httpModel.RegistrationStatus, error) {
	regs, err := u.registerAvailableRepository.RegistrationsAmount(ctx)
	if err != nil {
		return httpModel.RegistrationStatus(""), err
	}

	restrictions, err := u.registerAvailableRepository.SelectRegistrationRestrictions(ctx)
	if err != nil {
		return httpModel.RegistrationStatus(""), err
	}

	if regs < restrictions.Reserve {
		return httpModel.Available, nil
	}

	if regs < restrictions.Closed {
		return httpModel.Reserve, nil
	}

	return httpModel.Closed, nil
}
