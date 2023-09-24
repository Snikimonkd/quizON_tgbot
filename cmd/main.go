package main

import (
	"context"
	"encoding/json"
	"net/http"
	"quizon_bot/internal/app/delivery"
	"quizon_bot/internal/app/repository"
	"quizon_bot/internal/app/usecase"
	"quizon_bot/internal/config"
	"quizon_bot/internal/logger"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:     []string{"http://localhost:8000/front", "http://localhost:8000", "https://quiz-on.ru"},
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:     []string{"Content-Type", "Origin", "Accept", "Access-Control-Allow-Headers", "Access-Control-Expose-Headers", "Access-Control-Allow-Origin", "Authorization", "X-Requested-With", "X-CSRF-Token"},
		AllowCredentials:   true,
		OptionsPassthrough: true,
		MaxAge:             300,
	}))

	fs := http.FileServer(http.Dir("./front"))
	r.Handle("/front/*", http.StripPrefix("/front/", fs))

	r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
		type Kek struct {
			Lol string
		}
		kek := Kek{Lol: "hello"}
		b, err := json.Marshal(kek)
		err = json.NewEncoder(w).Encode(b)
		if err != nil {
			logger.Errorf("can't write body: %v", err)
		}
	})

	go func() {
		err := http.ListenAndServe(":8000", r)
		if err != nil {
			logger.Fatalf("can't start server: %v", err)
		}
		return
	}()

	ctx := context.Background()
	db := config.ConnectToPostgres(ctx)
	repository := repository.NewRepository(db)
	usecase := usecase.NewUsecase(repository)

	bot := config.ConnectToBot()

	app := delivery.NewBotDelivery(bot, usecase)
	app.ListenAndServe(ctx)
}
