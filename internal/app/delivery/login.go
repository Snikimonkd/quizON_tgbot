package delivery

import (
	"context"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type LoginUsecase interface {
	Login(ctx context.Context, userID int64, secretKey string) (bool, error)
}

// Login - получить игры
func (d *delivery) Login(ctx context.Context, update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	args := strings.Split(update.Message.CommandArguments(), " ")
	if len(args) < 1 {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "без пароля не прокатит"), nil
	}

	ok, err := d.loginUsecase.Login(ctx, update.Message.From.ID, args[0])
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, errorMessage), err
	}
	if !ok {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "неплохая попытка, попробуй еще раз"), nil
	}

	return tgbotapi.NewMessage(update.Message.Chat.ID, "welcome to the club buddy"), nil
}
