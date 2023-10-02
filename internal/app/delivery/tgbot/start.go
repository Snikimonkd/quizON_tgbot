package tgbot

import (
	"context"

	tgbotapi "github.com/matterbridge/telegram-bot-api/v6"
)

type StartUsecase interface {
	Start(ctx context.Context, userID int64) (tgbotapi.MessageConfig, error)
}

// Start - начало диалога
func (d *delivery) Start(ctx context.Context, update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	msg, err := d.startUsecase.Start(ctx, update.Message.From.ID)
	if err != nil {
		return msg, err
	}

	return msg, nil
}
