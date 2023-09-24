package delivery

import (
	"context"

	tgbotapi "github.com/matterbridge/telegram-bot-api/v6"
)

type RegisterStatesUsecase interface {
	RegisterStates(ctx context.Context, userID int64, msg string) (tgbotapi.MessageConfig, error)
}

func (d delivery) RegisterStates(ctx context.Context, update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	msg, err := d.registerStatesUsecase.RegisterStates(ctx, update.Message.From.ID, update.Message.Text)
	msg.ChatID = update.Message.From.ID
	return msg, err
}
