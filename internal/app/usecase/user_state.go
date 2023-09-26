package usecase

import (
	"context"
	"fmt"
	"quizon_bot/internal/generated/postgres/public/model"

	tgbotapi "github.com/matterbridge/telegram-bot-api/v6"
)

const regNo string = `😪 Вот теперь реально КвизOFF...

Будем очень ждать вашу команду на следующих играх! А чтобы точно успеть зарегистрироваться, следите за нашими социальными сетями и подписывайтесь на уведомления: 
📎 [группа ВКонтакте](https://vk.com/quizonmsk)
📎 [Telegram-канал](https://t.me/quizonmsk)

До встречи!`

const regSuccess string = `🤓 Поздравляем, первые задания от КвизON успешно выполнены — твоя команда зарегистрирована!

Посмотрим, как ты справишься с другими вопросами на первой игре новой Бауманской лиги КвизON. Напомним, что игра пройдет: 
⚡️4 октября, 19:00 
⚡️345 ауд. (ГУК) 

Если поменяются планы или возникнут другие вопросы, то пиши [Маше](https://t.me/maria_ilinyh).`

const maxTeamsAmountReached string = `Упс... 👉🏻👈🏻

К сожалению, места на игру 4 октября закончились.

Но не стоит опускать руки раньше времени! Мы можем зарегистрировать твою команду в резерв, чтобы в случае отказа от какой-то из прошедших команд, вы могли занять их место. 

Хочешь зарегистрироваться в резерв?`

const maxTeamsAmount int64 = 23

type State string

const (
	EMPTY           State = "empty"
	REG_IS_FULL     State = "reg_is_full"
	TEAM_ID         State = "team_id"
	TEAM_NAME       State = "team_name"
	CAPTAIN_NAME    State = "captain_name"
	GROUP_NAME      State = "group_name"
	PHONE           State = "phone"
	AMOUNT          State = "amount"
	QUIZON_QUESTION State = "quizon_question"
	REG_END         State = "reg_end"
	ONE_MORE_TEAM   State = "one_more_team"
)

type UserStatesHandlerRepository interface {
	UpdateState(ctx context.Context, state model.UserState) error
	GetState(ctx context.Context, userID int64) (string, error)
	GetRegistrationDraft(ctx context.Context, userID int64) (model.RegistrationsDraft, error)
	GenerateTeamID(ctx context.Context) (int64, error)
	UpdateRegistrationDraft(ctx context.Context, in model.RegistrationsDraft) error
	CreateRegistration(ctx context.Context, in model.Registrations) error
	CheckTeamsAmount(ctx context.Context) (int64, error)
}

func (u usecase) HandleUserState(ctx context.Context, update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	userID := update.Message.From.ID
	nickname := update.Message.From.UserName

	response := tgbotapi.MessageConfig{}
	response.Text = DefaultErrorMessage
	response.ChatID = userID
	response.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	response.ParseMode = "Markdown"
	response.DisableWebPagePreview = true

	state, err := u.registerStatesRepository.GetState(ctx, userID)
	if err != nil {
		return response, err
	}

	switch state {
	case string(EMPTY):
		amount, err := u.registerStatesRepository.CheckTeamsAmount(ctx)
		if err != nil {
			return response, err
		}

		if amount >= maxTeamsAmount {
			newState := model.UserState{
				UserID: userID,
				State:  string(REG_IS_FULL),
			}
			err := u.registerStatesRepository.UpdateState(ctx, newState)
			if err != nil {
				return response, err
			}

			yesBtn := tgbotapi.NewKeyboardButton("Да")
			noBnt := tgbotapi.NewKeyboardButton("Нет")
			btnRow := tgbotapi.NewKeyboardButtonRow(yesBtn, noBnt)
			keyboard := tgbotapi.NewReplyKeyboard(btnRow)
			response.ReplyMarkup = &keyboard
			response.Text = maxTeamsAmountReached
			return response, nil
		}

		newState := model.UserState{
			UserID: userID,
			State:  string(CAPTAIN_NAME),
		}
		err = u.registerStatesRepository.UpdateState(ctx, newState)
		if err != nil {
			return response, err
		}

		response.Text = "Как тебя зовут? (Пример: Иванов Иван Иванович)"
		return response, nil
	case string(REG_IS_FULL):
		if update.Message.Text == "Да" {
			newState := model.UserState{
				UserID: userID,
				State:  string(CAPTAIN_NAME),
			}
			err = u.registerStatesRepository.UpdateState(ctx, newState)
			if err != nil {
				return response, err
			}

			response.Text = "Как тебя зовут? (Пример: Николай Эрнестович Бауман)"
			return response, nil
		}

		if update.Message.Text == "Нет" {
			response.Text = regNo
			return response, nil
		}

		yesBtn := tgbotapi.NewKeyboardButton("Да")
		noBnt := tgbotapi.NewKeyboardButton("Нет")
		btnRow := tgbotapi.NewKeyboardButtonRow(yesBtn, noBnt)
		keyboard := tgbotapi.NewReplyKeyboard(btnRow)
		response.ReplyMarkup = &keyboard
		response.Text = maxTeamsAmountReached
		return response, nil
	case string(CAPTAIN_NAME):
		teamID, err := u.registerStatesRepository.GenerateTeamID(ctx)
		if err != nil {
			return response, err
		}

		now := u.clock.Now()
		draft := model.RegistrationsDraft{
			UserID:      userID,
			TgContact:   nickname,
			CreatedAt:   now,
			UpdatedAt:   now,
			TeamID:      teamID,
			CaptainName: &update.Message.Text,
		}

		err = u.registerStatesRepository.UpdateRegistrationDraft(ctx, draft)
		if err != nil {
			return response, err
		}

		newState := model.UserState{
			UserID: userID,
			State:  string(GROUP_NAME),
		}
		err = u.registerStatesRepository.UpdateState(ctx, newState)
		if err != nil {
			return response, err
		}

		response.Text = "Твоя учебная группа (Пример: СМ1-11)"
		return response, nil
	case string(GROUP_NAME):
		draft, err := u.registerStatesRepository.GetRegistrationDraft(ctx, userID)
		if err != nil {
			return response, err
		}

		draft.GroupName = &update.Message.Text
		err = u.registerStatesRepository.UpdateRegistrationDraft(ctx, draft)
		if err != nil {
			return response, err
		}

		newState := model.UserState{
			UserID: userID,
			State:  string(PHONE),
		}
		err = u.registerStatesRepository.UpdateState(ctx, newState)
		if err != nil {
			return response, err
		}

		response.Text = "Номер телефона (Пример: 8(999)888-77-66)"
		return response, nil
	case string(PHONE):
		draft, err := u.registerStatesRepository.GetRegistrationDraft(ctx, userID)
		if err != nil {
			return response, err
		}

		draft.Phone = &update.Message.Text
		err = u.registerStatesRepository.UpdateRegistrationDraft(ctx, draft)
		if err != nil {
			return response, err
		}

		newState := model.UserState{
			UserID: userID,
			State:  string(TEAM_NAME),
		}
		err = u.registerStatesRepository.UpdateState(ctx, newState)
		if err != nil {
			return response, err
		}

		response.Text = "Название вашей команды"
		return response, nil
	case string(TEAM_NAME):
		draft, err := u.registerStatesRepository.GetRegistrationDraft(ctx, userID)
		if err != nil {
			return response, err
		}

		draft.TeamName = &update.Message.Text
		err = u.registerStatesRepository.UpdateRegistrationDraft(ctx, draft)
		if err != nil {
			return response, err
		}

		newState := model.UserState{
			UserID: userID,
			State:  string(AMOUNT),
		}
		err = u.registerStatesRepository.UpdateState(ctx, newState)
		if err != nil {
			return response, err
		}

		response.Text = "Сколько человек в вашей команде? (В команде может быть от 4 до 8 человек)"
		return response, nil
	case string(AMOUNT):
		draft, err := u.registerStatesRepository.GetRegistrationDraft(ctx, userID)
		if err != nil {
			return response, err
		}

		draft.Amount = &update.Message.Text

		reg := model.Registrations{
			UserID:      draft.UserID,
			TgContact:   draft.TgContact,
			TeamID:      draft.TeamID,
			Pnohe:       *draft.Phone,
			TeamName:    *draft.TeamName,
			CaptainName: *draft.CaptainName,
			GroupName:   *draft.GroupName,
			Amount:      *draft.Amount,
			CreatedAt:   draft.CreatedAt,
			UpdatedAt:   draft.UpdatedAt,
		}
		err = u.registerStatesRepository.CreateRegistration(ctx, reg)
		if err != nil {
			return response, err
		}

		newState := model.UserState{
			UserID: userID,
			State:  string(QUIZON_QUESTION),
		}
		err = u.registerStatesRepository.UpdateState(ctx, newState)
		if err != nil {
			return response, err
		}

		b := tgbotapi.NewKeyboardButton("КвизON")
		r := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{b})
		response.Text = "КвизOFF?"
		response.ReplyMarkup = &r
		return response, nil
	case string(QUIZON_QUESTION):
		if update.Message.Text != "КвизON" {
			b := tgbotapi.NewKeyboardButton("КвизON")
			r := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{b})
			response.Text = "КвизOFF?"
			response.ReplyMarkup = &r
			return response, nil
		}

		newState := model.UserState{
			UserID: userID,
			State:  string(REG_END),
		}
		err = u.registerStatesRepository.UpdateState(ctx, newState)
		if err != nil {
			return response, err
		}

		clown := tgbotapi.NewKeyboardButton("🤡")
		row := tgbotapi.NewKeyboardButtonRow(clown)
		keyboard := tgbotapi.NewReplyKeyboard(row)
		response.ReplyMarkup = &keyboard
		response.Text = regSuccess
		return response, nil
	case string(REG_END):
		newState := model.UserState{
			UserID: userID,
			State:  string(ONE_MORE_TEAM),
		}
		err = u.registerStatesRepository.UpdateState(ctx, newState)
		if err != nil {
			return response, err
		}

		yesBtn := tgbotapi.NewKeyboardButton("Да")
		noBnt := tgbotapi.NewKeyboardButton("Нет")
		btnRow := tgbotapi.NewKeyboardButtonRow(yesBtn, noBnt)
		keyboard := tgbotapi.NewReplyKeyboard(btnRow)
		response.ReplyMarkup = &keyboard
		response.Text = "Хочешь зарегистрировать еще одну команду?"
		return response, nil
	case string(ONE_MORE_TEAM):
		if update.Message.Text == "Нет" {
			newState := model.UserState{
				UserID: userID,
				State:  string(REG_END),
			}
			err = u.registerStatesRepository.UpdateState(ctx, newState)
			if err != nil {
				return response, err
			}

			response.Text = "Ну не хочешь, как хочешь..."
			return response, nil
		}

		if update.Message.Text == "Да" {
			amount, err := u.registerStatesRepository.CheckTeamsAmount(ctx)
			if err != nil {
				return response, err
			}

			if amount >= maxTeamsAmount {
				newState := model.UserState{
					UserID: userID,
					State:  string(REG_IS_FULL),
				}
				err := u.registerStatesRepository.UpdateState(ctx, newState)
				if err != nil {
					return response, err
				}

				yesBtn := tgbotapi.NewKeyboardButton("Да")
				noBnt := tgbotapi.NewKeyboardButton("Нет")
				btnRow := tgbotapi.NewKeyboardButtonRow(yesBtn, noBnt)
				keyboard := tgbotapi.NewReplyKeyboard(btnRow)
				response.ReplyMarkup = &keyboard
				response.Text = maxTeamsAmountReached
				return response, nil
			}

			newState := model.UserState{
				UserID: userID,
				State:  string(CAPTAIN_NAME),
			}
			err = u.registerStatesRepository.UpdateState(ctx, newState)
			if err != nil {
				return response, err
			}

			response.Text = "Как тебя зовут? (Пример: Иванов Иван Иванович)"
			return response, nil
		}
	}

	return response, fmt.Errorf("unknown state: %v", state)
}
