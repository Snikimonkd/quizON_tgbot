package timesupport

import (
	"quizon_bot/internal/logger"
	"time"
)

var (
	LocMsk *time.Location
)

func init() {
	var err error
	LocMsk, err = time.LoadLocation("Europe/Moscow")
	if err != nil {
		logger.Fatalf("can't load location: %v", err.Error())
	}
}
