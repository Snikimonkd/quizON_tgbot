package tgbot

import (
	"context"

	tgbotapi "github.com/matterbridge/telegram-bot-api/v6"
)

type TableUsecase interface {
	Table(context.Context, tgbotapi.Update) (tgbotapi.MessageConfig, error)
}

type AuthUsecase interface {
	CheckAuth(ctx context.Context, userID int64) (bool, error)
}

func (d *delivery) Table(ctx context.Context, update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	ok, err := d.authUsecase.CheckAuth(ctx, update.Message.From.ID)
	if err != nil {
		return tgbotapi.NewMessage(update.Message.From.ID, "ойой"), err
	}

	if !ok {
		return tgbotapi.NewMessage(update.Message.From.ID, "надо залогиниться"), nil
	}

	msg, err := d.tableUsecase.Table(ctx, update)
	if err != nil {
		return msg, err
	}

	return msg, nil
}
