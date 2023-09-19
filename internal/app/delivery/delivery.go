package delivery

import (
	"context"
	"quizon_bot/internal/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Usecase interface {
	GamesUsecase
	CreateUsecase
	LoginUsecase
	RegisterUsecase
	ListUsecase
	RegisterStatesUsecase
}

type TgBotHandle func(ctx context.Context, update tgbotapi.Update) (tgbotapi.MessageConfig, error)

type delivery struct {
	bot *tgbotapi.BotAPI

	routes map[string]TgBotHandle

	gamesUsecase          Usecase
	createUsecase         Usecase
	loginUsecase          Usecase
	registerUsecase       Usecase
	listUsecase           Usecase
	registerStatesUsecase Usecase
}

func NewBotDelivery(bot *tgbotapi.BotAPI, usecases Usecase) delivery {
	return delivery{
		bot:                   bot,
		gamesUsecase:          usecases,
		createUsecase:         usecases,
		loginUsecase:          usecases,
		registerUsecase:       usecases,
		listUsecase:           usecases,
		registerStatesUsecase: usecases,
	}
}

var commands []tgbotapi.BotCommand = []tgbotapi.BotCommand{
	// user
	{
		Command:     "games",
		Description: "список ближайших игр.",
	},
	{
		Command:     "register",
		Description: "регистрация на игру.",
	},
}

const errorMessage string = "Ой, что-то пошло не так"

func (d *delivery) ListenAndServe(ctx context.Context) {
	d.routes = map[string]TgBotHandle{
		"games":    d.Games,
		"create":   d.Create,
		"login":    d.Login,
		"register": d.Register,
		"list":     d.List,
		"start":    d.Start,
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := d.bot.GetUpdatesChan(u)

	// Loop through each update.
	for update := range updates {
		if update.Message == nil ||
			update.Message.From.IsBot ||
			update.Message.Chat.IsGroup() ||
			update.Message.Chat.IsChannel() ||
			update.Message.Chat.IsSuperGroup() {
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
				d.Send(res)
			}
			// dialog
		} else {
			// user_id == chat_id -> user
			if update.Message.Chat.ID == update.Message.From.ID {
				res, err := d.RegisterStates(ctx, update)
				if err != nil {
					logger.Errorf("register state error: %w", err)
				}

				if res.Text != "" {
					d.Send(res)
				}
			} else {
				// chat -> do nothing
			}
		}
	}
}
