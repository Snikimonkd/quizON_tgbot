package delivery

import (
	"context"
	"fmt"

	tgbotapi "github.com/matterbridge/telegram-bot-api/v6"
)

// Start - начало диалога
func (d delivery) Start(ctx context.Context, update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	b := tgbotapi.NewInlineKeyboardButtonWebApp("Зарегестрироваться", tgbotapi.WebAppInfo{
		URL: "https://quiz-on.ru",
	})
	r := tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{b})

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.ReplyMarkup = &r
	d.bot.Send(msg)

	msg.Text = "Список доступных команд:\n"
	for i := 0; i < len(commands); i++ {
		msg.Text += fmt.Sprintf("/%v - %v\n", commands[i].Command, commands[i].Description)
	}

	return msg, nil
}
