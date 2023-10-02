package tgbot

import (
	"context"
	"quizon_bot/internal/logger"
	"sync"

	tgbotapi "github.com/matterbridge/telegram-bot-api/v6"
)

type Usecase interface {
	LoginUsecase
	UserStateHandlerUsecase
	TableUsecase
	AuthUsecase
	StartUsecase
}

type TgBotHandle func(ctx context.Context, update tgbotapi.Update) (tgbotapi.MessageConfig, error)

type delivery struct {
	bot *tgbotapi.BotAPI

	m sync.Mutex

	routes map[string]TgBotHandle

	loginUsecase          LoginUsecase
	registerStatesUsecase UserStateHandlerUsecase
	tableUsecase          TableUsecase
	authUsecase           AuthUsecase
	startUsecase          StartUsecase
}

func NewBotDelivery(bot *tgbotapi.BotAPI, usecases Usecase) delivery {
	return delivery{
		m:                     sync.Mutex{},
		bot:                   bot,
		loginUsecase:          usecases,
		registerStatesUsecase: usecases,
		tableUsecase:          usecases,
		authUsecase:           usecases,
		startUsecase:          usecases,
	}
}

const errorMessage string = "Ой, что-то пошло не так"

func (d *delivery) ListenAndServe(ctx context.Context) {
	d.routes = map[string]TgBotHandle{
		"login": d.Login,
		"start": d.Start,
		"table": d.Table,
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := d.bot.GetUpdatesChan(u)

	for update := range updates {
		if (update.Message != nil && update.Message.Text == "") ||
			(update.Message != nil && update.Message.Chat != nil &&
				(update.Message.Chat.IsGroup() ||
					update.Message.Chat.IsChannel() ||
					update.Message.Chat.IsSuperGroup() ||
					update.Message.Sticker != nil)) {
			continue
		}

		// command
		if update.Message.IsCommand() {
			handler, ok := d.routes[update.Message.Command()]
			if !ok {
				logger.Errorf("unknown route: %v", update.Message.Command())
				continue
			}

			res, err := handler(ctx, update)
			if err != nil {
				logger.Errorf("command error: %w", err)
			}

			if res.Text != "" {
				go d.Send(res)
			}
		} else {
			res, err := d.HandleUserState(ctx, update)
			if err != nil {
				logger.Errorf("register state error: %w", err)
			}

			if res.Text != "" {
				go d.Send(res)
			}
		}
	}
}
