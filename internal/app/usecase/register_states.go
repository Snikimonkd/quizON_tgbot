package usecase

import (
	"context"
	"errors"
	"fmt"
	"quizon_bot/internal/generated/postgres/public/model"

	tgbotapi "github.com/matterbridge/telegram-bot-api/v6"
)

type State int8

const (
	UNKNOWN State = iota
	TEAM_ID
	TEAM_NAME
	CAPTAIN_NAME
	GROUP_NAME
	AMOUNT
	PHONE
)

type RegisterStatesRepository interface {
	GetRegistrationDraft(ctx context.Context, userID int64) (model.RegistrationsDraft, error)
	GenerateTeamID(ctx context.Context) (int64, error)
	CheckGameID(ctx context.Context, gameID int64) (bool, error)
	UpdateRegistrationDraft(ctx context.Context, in model.RegistrationsDraft) error
	CreateRegistration(ctx context.Context, in model.Registrations) error
}

func (u usecase) RegisterStates(ctx context.Context, userID int64, msg string) (tgbotapi.MessageConfig, error) {
	response := tgbotapi.MessageConfig{}
	response.Text = DefaultErrorMessage

	draft, err := u.registerStatesRepository.GetRegistrationDraft(ctx, userID)
	if errors.Is(err, ErrNotFound) {
		return response, nil
	}
	if err != nil {
		return response, err
	}

	var state State
	if draft.CaptainName == nil {
		state = CAPTAIN_NAME
	} else if draft.GroupName == nil {
		state = GROUP_NAME
	} else if draft.Phone == nil {
		state = PHONE
	} else if draft.TeamName == nil {
		state = TEAM_NAME
	} else if draft.Amount == nil {
		state = AMOUNT
	}

	draft.UpdatedAt = u.clock.Now()

	switch state {
	//		return "Напиши id команды, если у твоей команды его нет, то напиши \"нет\" и я его сгенерирую", nil
	//	case TEAM_ID:
	//		var teamID int64
	//		flag := false
	//		if strings.ToLower(msg) == "нет" {
	//			flag = true
	//			teamID, err = u.registerStatesRepository.GenerateTeamID(ctx)
	//			if err != nil {
	//				return DefaultErrorMessage, err
	//			}
	//		} else {
	//			teamID, err = strconv.ParseInt(msg, 10, 64)
	//			if err != nil {
	//				return "Не удается распознать id команды", fmt.Errorf("can't parse teamID: %w", err)
	//			}
	//		}
	//
	//		draft.TeamID = &teamID
	//		err = u.registerStatesRepository.UpdateRegistrationDraft(ctx, draft)
	//		if err != nil {
	//			return DefaultErrorMessage, err
	//		}
	//		if errors.Is(err, ErrTeamIdIsUsed) {
	//			return "Команда с таким ID уже зарегестрирована на игру", nil
	//		}
	//
	//		if flag {
	//			return fmt.Sprintf("Id твоей команды: %d\n, сохрани его и используй при регистрации на следующие игры\nНапиши название команды", teamID), nil
	//		}
	case CAPTAIN_NAME:
		teamID, err := u.registerStatesRepository.GenerateTeamID(ctx)
		if err != nil {
			return response, err
		}

		draft.CaptainName = &msg
		draft.TeamID = teamID

		err = u.registerStatesRepository.UpdateRegistrationDraft(ctx, draft)
		if err != nil {
			return response, err
		}
		response.Text = "Твоя учебная группа (Пример: РК6-52)"
		return response, nil
	case GROUP_NAME:
		draft.GroupName = &msg
		err := u.registerStatesRepository.UpdateRegistrationDraft(ctx, draft)
		if err != nil {
			return response, err
		}
		response.Text = "Номер телефона (Пример: 89156567645)"
		return response, nil
	case PHONE:
		draft.Phone = &msg
		err := u.registerStatesRepository.UpdateRegistrationDraft(ctx, draft)
		if err != nil {
			return response, err
		}
		response.Text = "Название вашей команды"
		return response, nil
	case TEAM_NAME:
		draft.TeamName = &msg
		err := u.registerStatesRepository.UpdateRegistrationDraft(ctx, draft)
		if err != nil {
			return response, err
		}
		response.Text = "Сколько человек в вашей команде? (В команде может быть от 4 до 8 человек)"
		return response, nil
	case AMOUNT:
		err = u.registerStatesRepository.UpdateRegistrationDraft(ctx, draft)
		if err != nil {
			return response, err
		}

		draft.Amount = &msg

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

		b := tgbotapi.NewKeyboardButton("КвизON")
		r := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{b})

		reply := tgbotapi.MessageConfig{}
		reply.Text = "КвизOFF?"
		reply.ReplyMarkup = &r

		return reply, nil
	}

	return response, fmt.Errorf("unknown state: %v", state)
}
