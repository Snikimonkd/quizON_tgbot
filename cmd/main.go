package main

import (
	"context"
	"net/http"
	"quizon_bot/internal/config"
	"quizon_bot/internal/logger"

	httpDelivery "quizon_bot/internal/app/delivery/http"

	tgbotDelivery "quizon_bot/internal/app/delivery/tgbot"
	"quizon_bot/internal/app/repository"
	"quizon_bot/internal/app/usecase"
)

func main() {
	ctx := context.Background()

	router := config.NewRouter()

	db := config.ConnectToPostgres(ctx)
	repository := repository.NewRepository(db)
	usecase := usecase.NewUsecase(repository)
	httpDelivery := httpDelivery.NewDelivery(usecase)

	botAPI := config.ConnectToBot()
	tgbot := tgbotDelivery.NewBotDelivery(botAPI, usecase)

	go func() {
		f := func() {
			defer func() {
				r := recover()
				if r != nil {
					logger.Error("panic recovered: %v", r)
				}
			}()

			logger.Infof("bot is ready")
			tgbot.ListenAndServe(ctx)
		}

		for {
			f()
		}
	}()

	router.Post("/register", httpDelivery.Register)
	router.Get("/registrations", httpDelivery.Registrations)
	router.Get("/register-available", httpDelivery.RegisterAvailable)

	logger.Info("server started on port 8080")
	err := http.ListenAndServe(":8000", router)
	if err != nil {
		logger.Fatalf("can't start server: %v", err)
	}
}
