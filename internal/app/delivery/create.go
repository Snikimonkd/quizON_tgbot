package delivery

import (
	"context"
	"quizon_bot/internal/generated/postgres/public/model"
	"strings"

	tgbotapi "github.com/matterbridge/telegram-bot-api/v6"
)

type CreateUsecase interface {
	Create(ctx context.Context, req model.Games) error
	CheckAuth(ctx context.Context, userID int64) (bool, error)
}

func (d *delivery) Create(ctx context.Context, update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	ok, err := d.createUsecase.CheckAuth(ctx, update.Message.From.ID)
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, errorMessage), err
	}
	if !ok {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "надо залогиниться"), nil
	}

	args := strings.Split(update.Message.CommandArguments(), " ")
	game, err := mapArgsIntoGame(args)
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, errorMessage), err
	}

	err = d.createUsecase.Create(ctx, game)
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, errorMessage), err
	}

	return tgbotapi.NewMessage(update.Message.Chat.ID, "игра создана"), nil
}
