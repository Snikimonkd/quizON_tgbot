package http

import (
	"net/http"
	"quizon_bot/internal/logger"

	httpModel "quizon_bot/internal/app/delivery/http/model"

	"context"
)

type RegistrationsUsecase interface {
	Registrations(ctx context.Context) ([]httpModel.Registration, error)
}

func (d *delivery) Registrations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	res, err := d.registrationsUsecase.Registrations(ctx)
	if err != nil {
		logger.Error(err.Error())
		ResponseWithJson(w, http.StatusInternalServerError, Error{Msg: err.Error()})
		return
	}

	ResponseWithJson(w, http.StatusOK, res)

}
