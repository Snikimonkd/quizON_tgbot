package delivery

import (
	"context"
	"fmt"
	"quizon_bot/internal/generated/postgres/public/model"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"
)

type ListUsecase interface {
	List(ctx context.Context, gameID int) ([]model.Registrations, error)
	CheckAuth(ctx context.Context, userID int64) (bool, error)
}

func (d *delivery) List(ctx context.Context, update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	ok, err := d.createUsecase.CheckAuth(ctx, update.Message.From.ID)
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, errorMessage), err
	}
	if !ok {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "надо залогиниться"), nil
	}

	args := strings.Split(update.Message.CommandArguments(), " ")
	if len(args) < 1 {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "без пароля не прокатит"), nil
	}

	gameID, err := strconv.Atoi(args[0])
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, errorMessage), fmt.Errorf("can't atoi gameID: %w", err)
	}

	res, err := d.listUsecase.List(ctx, gameID)
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, errorMessage), err
	}

	ret := lo.Reduce(res, func(r string, t model.Registrations, i int) string {
		r += fmt.Sprintf("%d. Навзание: %v, ID: %v\n", i+1, t.TeamName, t.TeamID)
		return r
	}, "Список команд, зарегестрировавшихся на игру:\n")

	return tgbotapi.NewMessage(update.Message.Chat.ID, ret), nil
}
