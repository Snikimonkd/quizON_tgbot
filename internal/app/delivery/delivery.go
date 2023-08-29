package delivery

import (
	"context"
	"quizon_bot/internal/app/repository"
	"quizon_bot/internal/app/usecase"
	"quizon_bot/internal/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
)

type delivery struct {
	bot *tgbotapi.BotAPI

	gamesUsecase    GamesUsecase
	createUsecase   CreateUsecase
	loginUsecase    LoginUsecase
	registerUsecase RegisterUsecase
	listUsecase     ListUsecase
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

func NewBotDelivery(bot *tgbotapi.BotAPI, db *pgx.Conn) delivery {
	repositories := repository.NewRepository(db)
	usecases := usecase.NewUsecase(repositories)

	_, err := bot.Request(tgbotapi.NewDeleteMyCommands())
	if err != nil {
		logger.Fatalf("can't init bot: %w", err)
	}

	bot.Request(tgbotapi.NewDeleteMyCommands())
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

const errorMessage string = "Ой, что-то пошло не так"

func buildHelpMessage() string {
	res := "Команды для участников:\n"
	for i := 0; i < 2; i++ {
		res += "/" + commands[i].Command + " - " + commands[i].Description + "\n"
	}
	res += "\nКоманды для админов:\n"
	for i := 2; i < 5; i++ {
		res += "/" + commands[i].Command + " - " + commands[i].Description + "\n"
	}

	return res
}

func (d *delivery) ListenAndServe(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := d.bot.GetUpdatesChan(u)
	tgbotapi.NewBotCommandScopeDefault()

	// Loop through each update.
	for update := range updates {
		if update.Message == nil {
			continue
		}
		switch update.Message.Command() {
		case "games":
			msg, err := d.Games(ctx, update)
			if err != nil {
				logger.Errorf("Games() error = %w", err)
			}
			d.Send(msg)
		case "create":
			msg, err := d.Create(ctx, update)
			if err != nil {
				logger.Errorf("Create() error = %w", err)
			}
			d.Send(msg)
		case "login":
			msg, err := d.Login(ctx, update)
			if err != nil {
				logger.Errorf("Login() error = %w", err)
			}
			d.Send(msg)
		case "register":
			msg, err := d.Register(ctx, update)
			if err != nil {
				logger.Errorf("Register() error = %w", err)
			}
			d.Send(msg)
		case "list":
			msg, err := d.List(ctx, update)
			if err != nil {
				logger.Errorf("List() error = %w", err)
			}
			d.Send(msg)

		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, buildHelpMessage())
			d.bot.Send(msg)
		}
	}
}
