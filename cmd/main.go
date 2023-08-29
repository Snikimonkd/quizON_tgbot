package main

import (
	"context"
	"quizon_bot/internal/app/delivery"
	"quizon_bot/internal/config"
)

func main() {
	ctx := context.Background()
	db := config.ConnectToPostgres(ctx)
	bot := config.ConnectToBot()

	app := delivery.NewBotDelivery(bot, db)
	app.ListenAndServe(ctx)
}
