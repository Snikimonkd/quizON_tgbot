package delivery

import (
	"context"

	tgbotapi "github.com/matterbridge/telegram-bot-api/v6"
)

type TableUsecase interface {
	Table(context.Context, tgbotapi.Update) (tgbotapi.MessageConfig, error)
}

func (d delivery) Table(ctx context.Context, update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	msg, err := d.tableUsecase.Table(ctx, update)
	if err != nil {
		return msg, err
	}

	return msg, nil
}
