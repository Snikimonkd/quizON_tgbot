package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"quizon_bot/internal/logger"
)

type Usecase interface {
	RegisterUsecase
	RegistrationsUsecase
}

type delivery struct {
	registerUsecase      RegisterUsecase
	registrationsUsecase RegistrationsUsecase
}

func NewDelivery(usecase Usecase) *delivery {
	return &delivery{
		registerUsecase:      usecase,
		registrationsUsecase: usecase,
	}
}

type Error struct {
	Msg string `json:"msg"`
}

func ResponseWithJson(w http.ResponseWriter, code int, body interface{}) {
	w.WriteHeader(code)
	if body == nil {
		return
	}

	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		logger.Errorf("can't write body: %v", err)
	}

	return
}

func UnmarshalRequest(body io.ReadCloser, value interface{}) error {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(value)
	if err != nil {
		return fmt.Errorf("can't unmarshal request: %w", err)
	}

	return nil
}
