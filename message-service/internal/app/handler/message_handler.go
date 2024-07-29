package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/zer0day88/cme/message-service/helper"
	"github.com/zer0day88/cme/message-service/internal/app/model"
	"github.com/zer0day88/cme/message-service/internal/app/service"
	"github.com/zer0day88/cme/message-service/pkg/response"
)

type messageHandler struct {
	log        zerolog.Logger
	messageSrv *service.MessageService
}

type MessageHandler interface {
	Ok(c echo.Context) error
	Send(c echo.Context) error
	GetMessageHistory(c echo.Context) error
}

func NewUserHandler(log zerolog.Logger, messageSrv *service.MessageService) MessageHandler {

	return &messageHandler{log: log, messageSrv: messageSrv}
}

func (h *messageHandler) Ok(c echo.Context) error {
	return response.WriteJSON(c, response.OKNoError)
}

func (h *messageHandler) Send(c echo.Context) error {

	var req model.SendRequest
	userID := fmt.Sprint(c.Get("id"))

	if err := c.Bind(&req); err != nil {
		return response.WriteJSON(c, response.ErrBadRequest)
	}

	errv := helper.ValidateStruct(req)

	if errv != nil {
		return response.WriteJSON(c, response.ErrBadRequest.WithMsg(*errv))
	}

	resp := h.messageSrv.Send(c.Request().Context(), userID, req)

	if resp.Err != nil {
		h.log.Err(resp.Err).Send()
	}

	return response.WriteJSON(c, resp)

}

func (h *messageHandler) GetMessageHistory(c echo.Context) error {

	userID := fmt.Sprint(c.Get("id"))

	resp, api := h.messageSrv.GetByRecipient(c.Request().Context(), userID)

	if api.Err != nil {
		h.log.Err(api.Err).Send()
	}

	return response.WriteJSON(c, api.WithData(resp))
}
