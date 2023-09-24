package delivery

import (
	"context"
	"quizon_bot/internal/logger"

	tgbotapi "github.com/matterbridge/telegram-bot-api/v6"
)

type RegisterUsecase interface {
	Register(ctx context.Context, userID int64, nickname string) error
}

func (d delivery) Register(ctx context.Context, update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	err := d.registerUsecase.Register(ctx, update.Message.From.ID, update.Message.Chat.UserName)
	if err != nil {
		logger.Errorf("register error: %w", err)
		return tgbotapi.NewMessage(update.Message.From.ID, "kek"), nil
	}

	msg := "Как тебя зовут? (Пример: Иванов Иван Иванович)"
	return tgbotapi.NewMessage(update.Message.Chat.ID, msg), nil
}
