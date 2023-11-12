package utils

import (
	"time"

	"quizon_bot/internal/logger"
)

// LocMsk - default timezone
var LocMsk *time.Location

func init() {
	var err error
	LocMsk, err = time.LoadLocation("Europe/Moscow")
	if err != nil {
		logger.Fatalf("can't load Europe/Moscow location: %s", err)
	}
}

// PrettyTime - форматирует время в удобный для чтения формат
func PrettyTime(time time.Time) string {
	return time.In(LocMsk).Format("2006-01-02 15:04:05")
}
