package usecase

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

type StartRepository interface {
	Start(ctx context.Context, userID int64) error
}

func (u usecase) Start(ctx context.Context, userID int64) (tgbotapi.MessageConfig, error) {
	//	msg := tgbotapi.NewMessage(userID, DefaultErrorMessage)
	//	msg.ParseMode = "Markdown"
	//	msg.DisableWebPagePreview = true
	//
	//	err := u.startRepository.Start(ctx, userID)
	//	if err != nil {
	//		return msg, err
	//	}
	//
	//	msg.Text = startText
	//	btn := tgbotapi.NewKeyboardButton("Зарегистрироваться")
	//	row := tgbotapi.NewKeyboardButtonRow(btn)
	//	keyboard := tgbotapi.NewReplyKeyboard(row)
	//	msg.ReplyMarkup = &keyboard

	msg := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID: userID,
		},
		Text: "Жми на кнопку",
	}

	btn := tgbotapi.NewInlineKeyboardButtonWebApp("Зрегистрироваться", tgbotapi.WebAppInfo{URL: "https://quiz-on.ru"})
	row := tgbotapi.NewInlineKeyboardRow(btn)
	markup := tgbotapi.NewInlineKeyboardMarkup(row)
	msg.ReplyMarkup = &markup

	return msg, nil
}
