package tgbot

import (
	"quizon_bot/internal/logger"
	"time"

	tgbotapi "github.com/matterbridge/telegram-bot-api/v6"
)

// Send - отправляет сообщение с ретраями
func (d *delivery) Send(msg tgbotapi.MessageConfig) {
	for i := 0; i < 5; i++ {
		d.m.Lock()
		_, err := d.bot.Send(msg)
		d.m.Unlock()
		if err != nil {
			logger.Errorf("can't send message: %v", err)
			time.Sleep(time.Second * 3)
		} else {
			return
		}
	}
}
