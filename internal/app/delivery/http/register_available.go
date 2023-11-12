package http

import (
	"context"
	"net/http"

	httpModel "quizon_bot/internal/app/delivery/http/model"
	"quizon_bot/internal/logger"
)

type RegisterAvailableUsecase interface {
	RegisterAvailable(ctx context.Context) (httpModel.RegistrationStatus, error)
}

func (d *delivery) RegisterAvailable(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	status, err := d.registerAvailableUsecase.RegisterAvailable(ctx)
	if err != nil {
		logger.Error(err.Error())
		ResponseWithJson(w, http.StatusInternalServerError, Error{Msg: err.Error()})
		return
	}

	ResponseWithJson(w, http.StatusOK, httpModel.RegisterAvailable{Available: status})
}
