package main

import (
	"context"
	"net/http"

	httpDelivery "quizon_bot/internal/app/delivery/http"
	"quizon_bot/internal/app/repository"
	"quizon_bot/internal/app/usecase"
	"quizon_bot/internal/config"
	"quizon_bot/internal/logger"
)

func main() {
	ctx := context.Background()

	router := config.NewRouter()

	db := config.ConnectToPostgres(ctx)
	repository := repository.NewRepository(db)
	usecase := usecase.NewUsecase(repository)
	httpDelivery := httpDelivery.NewDelivery(usecase)

	router.Post("/register", httpDelivery.Register)
	router.Get("/registrations", httpDelivery.Registrations)
	router.Get("/register-available", httpDelivery.RegisterAvailable)
	//	router.Post("/login", httpDelivery.Login)

	logger.Info("server started on port 8080")
	err := http.ListenAndServe(":8000", router)
	if err != nil {
		logger.Fatalf("can't start server: %v", err)
	}
}
