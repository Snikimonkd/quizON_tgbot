package config

import (
	"quizon_bot/internal/logger"

	tgbotapi "github.com/matterbridge/telegram-bot-api/v6"
)

// ConnectToBot - подключается к боту
func ConnectToBot() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(GlobalConfig.Keys.Token)
	if err != nil {
		logger.Fatalf("can't init new bot api: %w", err)
	}

	bot.Debug = true

	logger.Infof("Authorized on account %s", bot.Self.UserName)

	return bot
}
