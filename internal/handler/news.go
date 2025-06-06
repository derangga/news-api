package handler

import (
	"newsapi/internal/exception"
	"newsapi/internal/model/request"
	responder "newsapi/internal/model/response"
	"newsapi/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type NewsHandler struct {
	uc        usecase.NewsUsecase
	validator *validator.Validate
}

func NewNewsArticlesHandler(
	validator *validator.Validate,
	uc usecase.NewsUsecase,
) NewsHandler {
	return NewsHandler{
		uc:        uc,
		validator: validator,
	}
}

func (h NewsHandler) CreateNews(c echo.Context) error {
	var req request.CreateNewsArticleRequest
	err := c.Bind(&req)
	if err != nil {
		log.Errorf("NewsHandler.bind: %w", err)
		return responder.ResponseBadRequest(c, "")
	}

	err = h.validator.Struct(req)
	if err != nil {
		log.Errorf("NewsHandler.validateStruct: %w", err)
		return responder.ResponseBadRequest(c, "create news article require [title, content, summary (optional), author_id, slug, status (optional), topicIDs]")
	}

	if err := h.uc.CreateNewsArticle(c.Request().Context(), req); err != nil {
		return responder.ResponseUnprocessableEntity(c, err.Error())
	}

	return responder.ResponseCreated(c, "news created")
}

func (h NewsHandler) GetNewsArticles(c echo.Context) error {
	articles, err := h.uc.GetNewsArticles(c.Request().Context())
	if err != nil {
		log.Errorf("NewsHandler.getTopics: %w", err)
		return responder.ResponseUnprocessableEntity(c, "failed get topics")
	}

	return responder.RespondOK(c, articles, "")
}

func (h NewsHandler) GetNewsBySlug(c echo.Context) error {
	slug := c.Param("slug")

	article, err := h.uc.GetNewsArticleBySlug(c.Request().Context(), slug)
	if err != nil {
		return responder.ResponseUnprocessableEntity(c, err.Error())
	}

	return responder.RespondOK(c, article, "")
}

func (h NewsHandler) UpdateNewsArticle(c echo.Context) error {
	slug := c.Param("slug")
	var req request.UpdateNewsArticleRequest
	err := c.Bind(&req)
	if err != nil {
		log.Errorf("NewsHandler.bind: %w", err)
		return responder.ResponseBadRequest(c, "")
	}

	err = h.validator.Struct(req)
	if err != nil {
		log.Errorf("NewsHandler.validateStruct: %w", err)
		return responder.ResponseBadRequest(c, "update news require one of [title, content, summary (optional), slug, status (optional), topicIDs]")
	}

	if err := h.uc.UpdateNewsArticleBySlug(c.Request().Context(), slug, req); err != nil {
		if customErr, ok := err.(exception.CustomError); ok && customErr.Code == 20005 {
			return responder.RespondOK(c, nil, "no field updated")
		}
		return responder.ResponseUnprocessableEntity(c, err.Error())
	}

	return responder.RespondOK(c, nil, "news updated")
}

func (h NewsHandler) DeleteNewsArticle(c echo.Context) error {
	slug := c.Param("slug")

	err := h.uc.DeleteNewsArticleBySlug(c.Request().Context(), slug)
	if err != nil {
		return responder.ResponseUnprocessableEntity(c, err.Error())
	}
	return responder.RespondOK(c, nil, "news deleted")
}
