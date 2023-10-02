package usecase

import (
	"context"
	"quizon_bot/internal/generated/postgres/public/model"

	tgbotapi "github.com/matterbridge/telegram-bot-api/v6"
)

func (u usecase) regIsFul(ctx context.Context, userID int64, text string) (tgbotapi.MessageConfig, error) {
	response := defaultMessage(userID)

	tx, err := u.registerStatesRepository.Begin(ctx)
	if err != nil {
		return response, err
	}

	if text == "Да" {
		newState := model.UserState{
			UserID: userID,
			State:  string(CAPTAIN_NAME),
		}
		err = u.registerStatesRepository.UpdateState(ctx, tx, newState)
		if err != nil {
			return response, err
		}

		err = u.registerStatesRepository.Commit(ctx, tx)
		if err != nil {
			return response, err
		}

		response.Text = "Как тебя зовут? (Пример: Николай Эрнестович Бауман)"
		return response, nil
	}

	if text == "Нет" {
		response.Text = regNo
		response.ReplyMarkup = CreateButtonKeyboard("Зарегистрироваться")
		return response, nil
	}

	response.ReplyMarkup = CreateButtonKeyboard("Зарегистрироваться")
	response.Text = maxTeamsAmountReached
	return response, nil

}
