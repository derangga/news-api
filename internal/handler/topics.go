package handler

import (
	"newsapi/internal/exception"
	"newsapi/internal/model/request"
	responder "newsapi/internal/model/response"
	"newsapi/internal/usecase"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type TopicsHandler struct {
	uc        usecase.TopicsUsecase
	validator *validator.Validate
}

func NewTopicsHandler(
	validator *validator.Validate,
	uc usecase.TopicsUsecase,
) TopicsHandler {
	return TopicsHandler{
		uc:        uc,
		validator: validator,
	}
}

func (h TopicsHandler) CreateTopic(c echo.Context) error {
	var req request.CreateTopicRequest
	err := c.Bind(&req)
	if err != nil {
		log.Errorf("TopicHandler.bind: %w", err)
		return responder.ResponseBadRequest(c, "")
	}

	err = h.validator.Struct(req)
	if err != nil {
		log.Errorf("TopicHandler.validateStruct: %w", err)
		return responder.ResponseBadRequest(c, "create topic require [name, description, slug]")
	}

	if err := h.uc.CreateTopic(c.Request().Context(), req); err != nil {
		return responder.ResponseUnprocessableEntity(c, err.Error())
	}

	return responder.ResponseCreated(c, "topic created")
}

func (h TopicsHandler) GetTopics(c echo.Context) error {
	topics, err := h.uc.GetTopics(c.Request().Context())
	if err != nil {
		log.Errorf("TopicHandler.getTopics: %w", err)
		return responder.ResponseUnprocessableEntity(c, "failed get topics")
	}

	return responder.RespondOK(c, topics, "")
}

func (h TopicsHandler) UpdateTopic(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return responder.ResponseBadRequest(c, "invalid id")
	}

	var req request.UpdateTopicRequest
	err = c.Bind(&req)
	if err != nil {
		log.Errorf("TopicHandler.bind: %w", err)
		return responder.ResponseBadRequest(c, "")
	}

	err = h.validator.Struct(req)
	if err != nil {
		log.Errorf("TopicHandler.validateStruct: %w", err)
		return responder.ResponseBadRequest(c, "update topic require one of [name, description, slug]")
	}

	req.ID = id
	if err := h.uc.UpdateTopic(c.Request().Context(), req); err != nil {
		if customErr, ok := err.(exception.CustomError); ok && customErr.Code == 20005 {
			return responder.RespondOK(c, nil, "no field updated")
		}
		return responder.ResponseUnprocessableEntity(c, err.Error())
	}

	return responder.RespondOK(c, nil, "topic updated")
}
