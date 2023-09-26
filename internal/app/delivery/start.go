package delivery

import (
	"context"

	tgbotapi "github.com/matterbridge/telegram-bot-api/v6"
)

const startText string = `КвизOFF? 

Если ты знаешь ответ на этот вопрос, то, скорее всего, ты уже с нами знаком. А если нет, то запоминай👇🏻

[КвизON](https://t.me/quizonmsk) — командная интеллектуально-развлекательная игра в формате викторины. Базируемся в МГТУ им. Н.Э. Баумана и устраиваем битвы логики и эрудиции среди студентов лучшего технического.

И ты попал в наш чат-бот, потому что захотел зарегистрироваться на ближайшую из игр: 
⚡️4 октября, 19:00 
⚡️345 ауд. (ГУК) 

Для регистрации жми кнопку *Зарегистрироваться*`

// Start - начало диалога
func (d delivery) Start(ctx context.Context, update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, startText)
	btn := tgbotapi.NewKeyboardButton("Зарегистрироваться")
	row := tgbotapi.NewKeyboardButtonRow(btn)
	keyboard := tgbotapi.NewReplyKeyboard(row)
	msg.ReplyMarkup = &keyboard
	msg.ParseMode = "Markdown"
	msg.DisableWebPagePreview = true

	return msg, nil
}
