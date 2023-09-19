package usecase

import (
	"context"
	"errors"
	"fmt"
	"quizon_bot/internal/generated/postgres/public/model"
	"strconv"
)

type State int8

const (
	UNKNOWN State = iota
	GAME_ID
	TEAM_ID
	TEAM_NAME
)

type RegisterStatesRepository interface {
	GetRegistrationDraft(ctx context.Context, userID int64) (model.RegistrationsDraft, error)
	GenerateTeamID(ctx context.Context) (int64, error)
	CheckGameID(ctx context.Context, gameID int64) (bool, error)
	UpdateRegistrationDraft(ctx context.Context, in model.RegistrationsDraft) error
	CreateRegistration(ctx context.Context, in model.Registrations) error
}

func (u usecase) RegisterStates(ctx context.Context, userID int64, msg string) (string, error) {
	draft, err := u.registerStatesRepository.GetRegistrationDraft(ctx, userID)
	if errors.Is(err, ErrNotFound) {
		return "", nil
	}
	if err != nil {
		return DefaultErrorMessage, err
	}

	var state State
	if draft.GameID == nil {
		state = GAME_ID
	} else if draft.TeamID == nil {
		state = TEAM_ID
	} else if draft.TeamName == nil {
		state = TEAM_NAME
	}

	draft.UpdatedAt = u.clock.Now()

	switch state {
	case GAME_ID:
		gameID, err := strconv.ParseInt(msg, 10, 64)
		if err != nil {
			return "Не могу понять, введи номер игры (число)", fmt.Errorf("can't parse gameID: %w", err)
		}
		ok, err := u.registerStatesRepository.CheckGameID(ctx, int64(gameID))
		if err != nil {
			return DefaultErrorMessage, err
		}
		if !ok {
			return "Игры с таким номером нет, посмотреть доступные игры можно с помощью команды /games", nil
		}

		teamID, err := u.registerStatesRepository.GenerateTeamID(ctx)
		if err != nil {
			return DefaultErrorMessage, err
		}

		draft.GameID = &gameID
		draft.TeamID = &teamID
		err = u.registerStatesRepository.UpdateRegistrationDraft(ctx, draft)
		if err != nil {
			return DefaultErrorMessage, err
		}

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
		return "Напиши навзание команды", nil
	case TEAM_NAME:
		draft.TeamName = &msg
		err := u.registerStatesRepository.UpdateRegistrationDraft(ctx, draft)
		if err != nil {
			return DefaultErrorMessage, err
		}

		reg := model.Registrations{
			GameID:    *draft.GameID,
			TeamID:    *draft.TeamID,
			TeamName:  *draft.TeamName,
			UserID:    draft.UserID,
			CreatedAt: draft.CreatedAt,
		}
		err = u.registerStatesRepository.CreateRegistration(ctx, reg)
		if err != nil {
			return DefaultErrorMessage, err
		}

		return fmt.Sprintf("Поздравляем вы зарегестрировались на игру\nНомер игры: %d\nID команды: %d\nНазвание команды: %s\n", reg.GameID, reg.TeamID, reg.TeamName), nil
	}

	return DefaultErrorMessage, fmt.Errorf("unknown state: %v", state)
}
