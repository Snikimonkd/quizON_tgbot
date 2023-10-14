package http

import (
	"net/http"

	httpModel "quizon_bot/internal/app/delivery/http/model"
	"quizon_bot/internal/logger"

	"context"
)

type RegisterAvailableUsecase interface {
	RegisterAvailable(ctx context.Context) (bool, error)
}

func (d *delivery) RegisterAvailable(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ok, err := d.registerAvailableUsecase.RegisterAvailable(ctx)
	if err != nil {
		logger.Error(err.Error())
		ResponseWithJson(w, http.StatusInternalServerError, Error{Msg: err.Error()})
		return
	}

	ResponseWithJson(w, http.StatusOK, httpModel.RegisterAvailable{Ok: ok})
}
