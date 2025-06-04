package handler

import (
	"net/http"
	"newsapi/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
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

func (h TopicsHandler) GetTopics(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
	})
}
