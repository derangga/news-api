package handler

import (
	"newsapi/internal/usecase"

	"github.com/go-playground/validator/v10"
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
