package usecase

import (
	"context"
	httpModel "quizon_bot/internal/app/delivery/http/model"
	"quizon_bot/internal/generated/postgres/public/model"

	"github.com/samber/lo"
)

type RegistrationsRepository interface {
	Registrations(ctx context.Context) ([]model.Registrations, error)
}

func (u usecase) Registrations(ctx context.Context) ([]httpModel.Registration, error) {
	res, err := u.registrationsRepository.Registrations(ctx)
	if err != nil {
		return nil, err
	}

	ret := lo.Map(res, func(item model.Registrations, index int) httpModel.Registration {
		return httpModel.Registration{
			Number:      int64(index),
			TgContact:   item.TgContact,
			TeamID:      item.TeamID,
			TeamName:    item.TeamName,
			CaptainName: item.CaptainName,
			Phone:       item.Phone,
			GroupName:   item.GroupName,
			Amount:      item.Amount,
		}
	})

	return ret, nil
}