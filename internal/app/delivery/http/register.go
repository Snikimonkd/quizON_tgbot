package http

import (
	"math/rand"
	"net/http"
	httpModel "quizon_bot/internal/app/delivery/http/model"
	"quizon_bot/internal/logger"

	"context"
)

type RegisterUsecase interface {
	Register(ctx context.Context, req httpModel.Register) error
}

func (d *delivery) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req httpModel.Register
	err := UnmarshalRequest(r.Body, &req)
	if err != nil {
		logger.Error(err.Error())
		ResponseWithJson(w, http.StatusBadRequest, Error{Msg: err.Error()})
	}

	err = d.registerUsecase.Register(ctx, req)
	if err != nil {
		logger.Error(err.Error())
		ResponseWithJson(w, http.StatusBadRequest, Error{Msg: err.Error()})
	}

	_ = r.Context()

	if rand.Int63()%10 > 5 {
	} else {
		w.WriteHeader(500)
	}
}
