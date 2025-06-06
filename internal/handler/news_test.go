package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"newsapi/internal/exception"
	"newsapi/internal/handler"
	"newsapi/internal/model/response"
	mock_usecase "newsapi/mocks/usecase"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type NewsHandlerAccessor struct {
	newsUC  *mock_usecase.MockNewsUsecase
	handler handler.NewsHandler
}

func newNewsHandlerAccessor(ctrl *gomock.Controller) NewsHandlerAccessor {
	newsUC := mock_usecase.NewMockNewsUsecase(ctrl)
	validator := validator.New()
	handler := handler.NewNewsArticlesHandler(validator, newsUC)
	return NewsHandlerAccessor{
		newsUC:  newsUC,
		handler: handler,
	}
}

func Test_CreateNews(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	accessor := newNewsHandlerAccessor(ctrl)
	h := accessor.handler
	e := echo.New()

	validBody := `{
		"title": "A Valid Title",
		"content": "This is the valid content for a news article.",
		"summary": "A valid summary",
		"author_id": 1,
		"slug": "a-valid-title",
		"status": "draft",
		"topic_ids": [1, 2]
	}`

	tests := []struct {
		name      string
		body      string
		initMock  func()
		assertion func(echo.Context, *httptest.ResponseRecorder, error)
	}{
		{
			name:     "invalid payload returns 400",
			body:     `{}`, // missing required fields
			initMock: func() {},
			assertion: func(ctx echo.Context, rr *httptest.ResponseRecorder, err error) {
				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, rr.Code)
			},
		},
		{
			name: "valid payload but usecase returns error, expect 422",
			body: validBody,
			initMock: func() {
				accessor.newsUC.EXPECT().
					CreateNewsArticle(gomock.Any(), gomock.Any()).
					Return(errors.New("something went wrong"))
			},
			assertion: func(ctx echo.Context, rr *httptest.ResponseRecorder, err error) {
				assert.NoError(t, err)
				assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
			},
		},
		{
			name: "valid payload and successful creation, expect 201",
			body: validBody,
			initMock: func() {
				accessor.newsUC.EXPECT().
					CreateNewsArticle(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			assertion: func(ctx echo.Context, rr *httptest.ResponseRecorder, err error) {
				assert.NoError(t, err)
				assert.Equal(t, http.StatusCreated, rr.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.initMock()

			req := httptest.NewRequest(http.MethodPost, "/api/v1/news", strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := h.CreateNews(c)
			tt.assertion(c, rec, err)
		})
	}
}

func Test_GetNewsArticles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	accessor := newNewsHandlerAccessor(ctrl)
	h := accessor.handler
	e := echo.New()

	tests := []struct {
		name      string
		initMock  func()
		assertion func(echo.Context, *httptest.ResponseRecorder, error)
	}{
		{
			name: "usecase returns error, expect 422",
			initMock: func() {
				accessor.newsUC.EXPECT().
					GetNewsArticles(gomock.Any()).
					Return(nil, errors.New("unexpected error"))
			},
			assertion: func(ctx echo.Context, rr *httptest.ResponseRecorder, err error) {
				assert.NoError(t, err)
				assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
			},
		},
		{
			name: "usecase returns articles successfully, expect 200",
			initMock: func() {
				accessor.newsUC.EXPECT().
					GetNewsArticles(gomock.Any()).
					Return([]response.NewsArticle{}, nil)
			},
			assertion: func(ctx echo.Context, rr *httptest.ResponseRecorder, err error) {
				assert.NoError(t, err)
				assert.Equal(t, http.StatusOK, rr.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.initMock()

			req := httptest.NewRequest(http.MethodGet, "/api/v1/news", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := h.GetNewsArticles(c)
			tt.assertion(c, rec, err)
		})
	}
}

func Test_GetNewsBySlug(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	accessor := newNewsHandlerAccessor(ctrl)
	h := accessor.handler
	e := echo.New()

	tests := []struct {
		name      string
		slug      string
		initMock  func()
		assertion func(*httptest.ResponseRecorder, error)
	}{
		{
			name: "returns article successfully",
			slug: "test-slug",
			initMock: func() {
				accessor.newsUC.EXPECT().
					GetNewsArticleBySlug(gomock.Any(), "test-slug").
					Return(response.NewsArticleWithTopic{
						ID:    1,
						Title: "Test article",
					}, nil)
			},
			assertion: func(rr *httptest.ResponseRecorder, err error) {
				assert.NoError(t, err)
				assert.Equal(t, http.StatusOK, rr.Code)
			},
		},
		{
			name: "returns 422 on error",
			slug: "missing-slug",
			initMock: func() {
				accessor.newsUC.EXPECT().
					GetNewsArticleBySlug(gomock.Any(), "missing-slug").
					Return(response.NewsArticleWithTopic{}, errors.New("not found"))
			},
			assertion: func(rr *httptest.ResponseRecorder, err error) {
				assert.NoError(t, err)
				assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.initMock()

			req := httptest.NewRequest(http.MethodGet, "/api/v1/news/"+tt.slug, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("slug")
			c.SetParamValues(tt.slug)

			err := h.GetNewsBySlug(c)
			tt.assertion(rec, err)
		})
	}
}

func Test_UpdateNewsArticle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	accessor := newNewsHandlerAccessor(ctrl)
	h := accessor.handler
	e := echo.New()

	validBody := `{"title": "Updated Title"}`
	slug := "some-article"

	tests := []struct {
		name      string
		body      string
		initMock  func()
		assertion func(echo.Context, *httptest.ResponseRecorder, error)
	}{
		{
			name:     "invalid JSON body, expect 400",
			body:     `{"title": 123}`, // invalid type
			initMock: func() {},
			assertion: func(ctx echo.Context, rr *httptest.ResponseRecorder, err error) {
				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, rr.Code)
			},
		},
		{
			name:     "validation fails, expect 400",
			body:     `{"title": "abc"}`, // too short
			initMock: func() {},
			assertion: func(ctx echo.Context, rr *httptest.ResponseRecorder, err error) {
				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, rr.Code)
			},
		},
		{
			name: "usecase returns custom no update error, expect 200 with message",
			body: validBody,
			initMock: func() {
				accessor.newsUC.EXPECT().
					UpdateNewsArticleBySlug(gomock.Any(), slug, gomock.Any()).
					Return(exception.CustomError{Code: 20005})
			},
			assertion: func(ctx echo.Context, rr *httptest.ResponseRecorder, err error) {
				assert.NoError(t, err)
				assert.Equal(t, http.StatusOK, rr.Code)
			},
		},
		{
			name: "usecase returns generic error, expect 422",
			body: validBody,
			initMock: func() {
				accessor.newsUC.EXPECT().
					UpdateNewsArticleBySlug(gomock.Any(), slug, gomock.Any()).
					Return(errors.New("unexpected error"))
			},
			assertion: func(ctx echo.Context, rr *httptest.ResponseRecorder, err error) {
				assert.NoError(t, err)
				assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
			},
		},
		{
			name: "successfully updated, expect 200",
			body: validBody,
			initMock: func() {
				accessor.newsUC.EXPECT().
					UpdateNewsArticleBySlug(gomock.Any(), slug, gomock.Any()).
					Return(nil)
			},
			assertion: func(ctx echo.Context, rr *httptest.ResponseRecorder, err error) {
				assert.NoError(t, err)
				assert.Equal(t, http.StatusOK, rr.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.initMock()

			req := httptest.NewRequest(http.MethodPatch, "/api/v1/news/"+slug, strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("slug")
			c.SetParamValues(slug)

			err := h.UpdateNewsArticle(c)
			tt.assertion(c, rec, err)
		})
	}
}

func Test_DeleteNewsArticle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	accessor := newNewsHandlerAccessor(ctrl)
	h := accessor.handler
	e := echo.New()

	tests := []struct {
		name      string
		slug      string
		initMock  func()
		assertion func(*httptest.ResponseRecorder, error)
	}{
		{
			name: "deletes article successfully",
			slug: "test-slug",
			initMock: func() {
				accessor.newsUC.EXPECT().
					DeleteNewsArticleBySlug(gomock.Any(), "test-slug").
					Return(nil)
			},
			assertion: func(rr *httptest.ResponseRecorder, err error) {
				assert.NoError(t, err)
				assert.Equal(t, http.StatusOK, rr.Code)
			},
		},
		{
			name: "returns 422 on delete error",
			slug: "invalid-slug",
			initMock: func() {
				accessor.newsUC.EXPECT().
					DeleteNewsArticleBySlug(gomock.Any(), "invalid-slug").
					Return(errors.New("delete error"))
			},
			assertion: func(rr *httptest.ResponseRecorder, err error) {
				assert.NoError(t, err)
				assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.initMock()

			req := httptest.NewRequest(http.MethodDelete, "/api/v1/news/"+tt.slug, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("slug")
			c.SetParamValues(tt.slug)

			err := h.DeleteNewsArticle(c)
			tt.assertion(rec, err)
		})
	}
}
