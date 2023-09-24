package delivery

import (
	"context"
	"fmt"

	tgbotapi "github.com/matterbridge/telegram-bot-api/v6"
)

// Start - начало диалога
func (d delivery) Start(ctx context.Context, update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	msg := "Список доступных команд:\n"
	for i := 0; i < len(commands); i++ {
		msg += fmt.Sprintf("/%v - %v\n", commands[i].Command, commands[i].Description)
	}

	return tgbotapi.NewMessage(update.Message.Chat.ID, msg), nil
}
