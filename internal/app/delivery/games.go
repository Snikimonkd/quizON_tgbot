package delivery

import (
	"context"
	"quizon_bot/internal/generated/postgres/public/model"

	tgbotapi "github.com/matterbridge/telegram-bot-api/v6"
	"github.com/samber/lo"
)

// GamesUsecase - интерфейс для получения игр
type GamesUsecase interface {
	Games(ctx context.Context) ([]model.Games, error)
}

const NoGames string = "К сожалению на ближайшее время игр не запланировано 😢"

// Games - получить игры
func (d *delivery) Games(ctx context.Context, update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	games, err := d.gamesUsecase.Games(ctx)
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, errorMessage), err
	}
	if len(games) == 0 {
		return tgbotapi.NewMessage(update.Message.Chat.ID, NoGames), nil
	}

	res := lo.Map(games, func(t model.Games, i int) string {
		return mapGame(t)
	})
	st := lo.Reduce(res, func(r string, t string, i int) string {
		r += t + "---------------------------------\n"
		return r
	}, "")

	return tgbotapi.NewMessage(update.Message.Chat.ID, st), nil
}
