package main

import (
	"context"
	"quizon_bot/internal/app/delivery"
	"quizon_bot/internal/app/repository"
	"quizon_bot/internal/app/usecase"
	"quizon_bot/internal/config"
)

func main() {
	ctx := context.Background()
	db := config.ConnectToPostgres(ctx)
	repository := repository.NewRepository(db)
	usecase := usecase.NewUsecase(repository)

	bot := config.ConnectToBot()

	app := delivery.NewBotDelivery(bot, usecase)
	app.ListenAndServe(ctx)
}
