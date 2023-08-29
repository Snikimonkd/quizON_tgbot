package delivery

import (
	"context"
	"quizon_bot/internal/generated/postgres/public/model"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type RegisterUsecase interface {
	Register(ctx context.Context, req model.Registrations) error
}

func (d delivery) Register(ctx context.Context, update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	args := strings.Split(update.Message.CommandArguments(), " ")
	req, err := mapArgsIntoRegistration(args)
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "какие-то непонятные данные"), nil
	}

	err = d.registerUsecase.Register(ctx, req)
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, errorMessage), err
	}

	return tgbotapi.NewMessage(update.Message.Chat.ID, "поздравляю вы успешно зарегестрировались на игру"), nil
}
