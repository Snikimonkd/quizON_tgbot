package delivery

import (
	"context"

	tgbotapi "github.com/matterbridge/telegram-bot-api/v6"
)

type UserStateHandlerUsecase interface {
	HandleUserState(ctx context.Context, update tgbotapi.Update) (tgbotapi.MessageConfig, error)
}

func (d delivery) HandleUserState(ctx context.Context, update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	msg, err := d.registerStatesUsecase.HandleUserState(ctx, update)
	return msg, err
}
