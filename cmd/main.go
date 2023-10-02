package main

import (
	"context"
	"net/http"
	"quizon_bot/internal/config"
	"quizon_bot/internal/logger"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	delivery "quizon_bot/internal/app/delivery/http"
	"quizon_bot/internal/app/repository"
	"quizon_bot/internal/app/usecase"
)

func main() {
	ctx := context.Background()

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"localhost:3000", "http://localhost:3000", "https://quiz-on.ru", "https://www.quiz-on.ru"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Origin", "User-Agent", "Sec-Fetch-Site", "Sec-Fetch-Mode", "Sec-Fetch-Dest", "Referer", "Access-Control-Request-Method", "Access-Control-Request-Headers", "Accept-Language", "Accept-Encoding", "Accept", "Access-Control-Allow-Headers", "Access-Control-Expose-Headers", "Access-Control-Allow-Origin", "Authorization", "X-Requested-With", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
		Debug:            true,
	}))

	db := config.ConnectToPostgres(ctx)
	repository := repository.NewRepository(db)
	usecase := usecase.NewUsecase(repository)
	delivery := delivery.NewDelivery(usecase)

	router.Post("/register", delivery.Register)

	logger.Info("server started on port 8080")
	err := http.ListenAndServe(":8000", router)
	if err != nil {
		logger.Fatalf("can't start server: %v", err)
	}

	// ctx := context.Background()
	// db := config.ConnectToPostgres(ctx)
	// repository := repository.NewRepository(db)
	// usecase := usecase.NewUsecase(repository)
	//
	// bot := config.ConnectToBot()
	//
	// app := delivery.NewBotDelivery(bot, usecase)
	// app.ListenAndServe(ctx)
}
