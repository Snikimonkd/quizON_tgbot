package delivery

import (
	"context"
	"quizon_bot/internal/logger"

	tgbotapi "github.com/matterbridge/telegram-bot-api/v6"
)

type RegisterUsecase interface {
	Register(ctx context.Context, userID int64) error
}

func (d delivery) Register(ctx context.Context, update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	err := d.registerUsecase.Register(ctx, update.Message.From.ID)
	if err != nil {
		logger.Errorf("register error: %w", err)
		return tgbotapi.NewMessage(update.Message.From.ID, "kek"), nil
	}

	msg := "Напиши номер игры на которую ты хочешь зарегестрироваться"
	return tgbotapi.NewMessage(update.Message.Chat.ID, msg), nil
}
