package handler

import (
	"newsapi/internal/model/request"
	responder "newsapi/internal/model/response"
	"newsapi/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type UsersHandler struct {
	uc        usecase.UsersUsecase
	validator *validator.Validate
}

func NewUsersHandler(
	validator *validator.Validate,
	uc usecase.UsersUsecase,
) UsersHandler {
	return UsersHandler{
		uc:        uc,
		validator: validator,
	}
}

func (h UsersHandler) CreateUser(c echo.Context) error {
	var req request.CreateUserRequest
	err := c.Bind(&req)
	if err != nil {
		log.Errorf("TopicHandler.bind: %w", err)
		return responder.ResponseBadRequest(c, "")
	}

	err = h.validator.Struct(req)
	if err != nil {
		log.Errorf("TopicHandler.validateStruct: %w", err)
		return responder.ResponseBadRequest(c, "create topic require [name, email]")
	}

	if err := h.uc.CreateUser(c.Request().Context(), req); err != nil {
		return responder.ResponseUnprocessableEntity(c, err.Error())
	}

	return responder.ResponseCreated(c, "user registered")
}
