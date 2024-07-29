package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/zer0day88/cme/user-service/helper"
	"github.com/zer0day88/cme/user-service/internal/app/model"
	"github.com/zer0day88/cme/user-service/internal/app/service"
	"github.com/zer0day88/cme/user-service/pkg/response"
)

type userHandler struct {
	log     zerolog.Logger
	userSrv *service.UserService
}

type UserHandler interface {
	Ok(c echo.Context) error
	Register(c echo.Context) error
	Login(c echo.Context) error
}

func NewUserHandler(log zerolog.Logger, userSrv *service.UserService) UserHandler {

	return &userHandler{log: log, userSrv: userSrv}
}

func (h *userHandler) Ok(c echo.Context) error {

	resp, api := h.userSrv.Get(c.Request().Context())
	return response.WriteJSON(c, api.WithData(resp))
}

func (h *userHandler) Register(c echo.Context) error {

	var req model.RegisterRequest

	if err := c.Bind(&req); err != nil {
		return response.WriteJSON(c, response.ErrBadRequest)
	}

	errv := helper.ValidateStruct(req)

	if errv != nil {
		return response.WriteJSON(c, response.ErrBadRequest.WithMsg(*errv))
	}

	resp := h.userSrv.Register(c.Request().Context(), req)

	if resp.Err != nil {
		h.log.Err(resp.Err).Send()
	}

	return response.WriteJSON(c, resp)

}

func (h *userHandler) Login(c echo.Context) error {
	var req model.LoginRequest

	if err := c.Bind(&req); err != nil {
		return response.WriteJSON(c, response.ErrBadRequest)
	}

	resp, api := h.userSrv.Login(c.Request().Context(), req)

	if api.Err != nil {
		h.log.Err(api.Err).Send()
	}

	return response.WriteJSON(c, api.WithData(resp))
}
