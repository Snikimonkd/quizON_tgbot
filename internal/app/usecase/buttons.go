package usecase

import tgbotapi "github.com/matterbridge/telegram-bot-api/v6"

func CreateYesNoKeyboard() *tgbotapi.ReplyKeyboardMarkup {
	yesBtn := tgbotapi.NewKeyboardButton("Да")
	noBnt := tgbotapi.NewKeyboardButton("Нет")
	btnRow := tgbotapi.NewKeyboardButtonRow(yesBtn, noBnt)
	keyboard := tgbotapi.NewReplyKeyboard(btnRow)
	return &keyboard
}

func CreateButtonKeyboard(str string)*tgbotapi.ReplyKeyboardMarkup {
	btn := tgbotapi.NewKeyboardButton(str)
	row := tgbotapi.NewKeyboardButtonRow(btn)
	keyboard := tgbotapi.NewReplyKeyboard(row)
	return &keyboard
}
