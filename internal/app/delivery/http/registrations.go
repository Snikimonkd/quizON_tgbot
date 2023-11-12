package http

import (
	"context"
	"net/http"

	httpModel "quizon_bot/internal/app/delivery/http/model"
	"quizon_bot/internal/logger"
)

type RegistrationsUsecase interface {
	Registrations(ctx context.Context) ([]httpModel.Registration, error)
}

func (d *delivery) Registrations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req httpModel.Registrations
	err := UnmarshalRequest(r.Body, &req)
	if err != nil {
		logger.Error(err.Error())
		ResponseWithJson(w, http.StatusBadRequest, Error{Msg: err.Error()})
		return
	}

	if req.Password != "09154cb6-f723-4f3d-943c-7a6e4b155eb1" {
		logger.Infof("registrations wit hwrong password: %v", req.Password)
		ResponseWithJson(w, http.StatusUnauthorized, Error{Msg: "ti po moemu chto-to pereputal"})
		return
	}

	res, err := d.registrationsUsecase.Registrations(ctx)
	if err != nil {
		logger.Error(err.Error())
		ResponseWithJson(w, http.StatusInternalServerError, Error{Msg: err.Error()})
		return
	}

	ResponseWithJson(w, http.StatusOK, res)
}
