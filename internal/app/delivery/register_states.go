package delivery

import (
	"context"

	tgbotapi "github.com/matterbridge/telegram-bot-api/v6"
)

type RegisterStatesUsecase interface {
	RegisterStates(ctx context.Context, userID int64, msg string) (string, error)
}

func (d delivery) RegisterStates(ctx context.Context, update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	msg, err := d.registerStatesUsecase.RegisterStates(ctx, update.Message.From.ID, update.Message.Text)
	return tgbotapi.NewMessage(update.Message.Chat.ID, msg), err
}
