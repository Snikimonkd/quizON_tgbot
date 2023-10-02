package usecase

import (
	"context"
	"quizon_bot/internal/generated/postgres/public/model"

	tgbotapi "github.com/matterbridge/telegram-bot-api/v6"
)

func (u usecase) empty(ctx context.Context, userID int64) (tgbotapi.MessageConfig, error) {
	response := defaultMessage(userID)

	tx, err := u.registerStatesRepository.Begin(ctx)
	if err != nil {
		return response, err
	}

	// проверяем количество уже зарегестрированных команд
	amount, err := u.registerStatesRepository.CheckTeamsAmount(ctx, tx)
	if err != nil {
		return response, err
	}

	// если текущее колиечество зарегестрированных команд >= максимального,
	// то предлагаем регу в резерв
	if amount >= maxTeamsAmount {
		newState := model.UserState{
			UserID: userID,
			State:  string(REG_IS_FULL),
		}

		err := u.registerStatesRepository.UpdateState(ctx, tx, newState)
		if err != nil {
			return response, err
		}

		err = u.registerStatesRepository.Commit(ctx, tx)
		if err != nil {
			return response, err
		}

		response.ReplyMarkup = CreateYesNoKeyboard()
		response.Text = maxTeamsAmountReached
		return response, nil
	}

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

	response.Text = "Как тебя зовут? (Пример: Иванов Иван Иванович)"
	return response, nil
}
