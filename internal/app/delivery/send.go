package delivery

import (
	"quizon_bot/internal/logger"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Send - отправляет сообщение с ретраями
func (d *delivery) Send(msg tgbotapi.MessageConfig) {
	for i := 0; i < 5; i++ {
		_, err := d.bot.Send(msg)
		if err != nil {
			logger.Errorf("can't send message: %v", err)
			time.Sleep(time.Minute)
		} else {
			return
		}
	}
}
