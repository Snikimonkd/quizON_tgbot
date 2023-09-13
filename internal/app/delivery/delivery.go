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
}

type TgBotHandle func(ctx context.Context, update tgbotapi.Update) (tgbotapi.MessageConfig, error)

type delivery struct {
	bot *tgbotapi.BotAPI

	routes map[string]TgBotHandle

	gamesUsecase    Usecase
	createUsecase   Usecase
	loginUsecase    Usecase
	registerUsecase Usecase
	listUsecase     Usecase
}

func NewBotDelivery(bot *tgbotapi.BotAPI, usecases Usecase) delivery {
	_, err := bot.Request(tgbotapi.NewDeleteMyCommands())
	if err != nil {
		logger.Fatalf("can't init bot: %w")
	}
	_, err = bot.Request(tgbotapi.NewSetMyCommands(commands...))
	if err != nil {
		logger.Fatalf("can't init bot: %w", err)
	}

	return delivery{
		bot:             bot,
		gamesUsecase:    usecases,
		createUsecase:   usecases,
		loginUsecase:    usecases,
		registerUsecase: usecases,
		listUsecase:     usecases,
	}
}

var commands []tgbotapi.BotCommand = []tgbotapi.BotCommand{
	{
		Command:     "games",
		Description: "список ближайших игр",
	},
	{
		Command:     "register",
		Description: "регистрация на игру. Пример: /register <№ игры> <id команды> <навзание команды>",
	},
	{
		Command:     "login",
		Description: "залогиниться как адми. Пример: /login <ключик>",
	},
	{
		Command:     "create",
		Description: "создать игру. Пример: /create <число> <месяц> <год> <часы:минуты> <место проведения> <опциональный комментарий>",
	},
	{
		Command:     "list",
		Description: "список команд, зарегестрировавшихся на игру. Пример: /list <№ игры>",
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
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := d.bot.GetUpdatesChan(u)

	// Loop through each update.
	for update := range updates {
		if update.Message == nil {
			continue
		}

		v, ok := d.routes[update.Message.Command()]
		if !ok {
			logger.Errorf("unknown route: %v", update.Message.Command())
		}

		res, err := v(ctx, update)
		if err != nil {
			logger.Errorf("error")
		}

		if res.Text != "" {
			d.Send(res)
		}
	}
}
