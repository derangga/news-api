package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"newsapi/internal/exception"
	"newsapi/internal/handler"
	"newsapi/internal/model/request"
	"newsapi/internal/model/response"
	"newsapi/internal/utils"
	mock_usecase "newsapi/mocks/usecase"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type TopicsHandlerAccessor struct {
	topicsUC *mock_usecase.MockTopicsUsecase
	handler  handler.TopicsHandler
}

func newTopicsHandlerAccessor(ctrl *gomock.Controller) TopicsHandlerAccessor {
	topicsUC := mock_usecase.NewMockTopicsUsecase(ctrl)
	validator := validator.New()
	handler := handler.NewTopicsHandler(validator, topicsUC)
	return TopicsHandlerAccessor{
		topicsUC: topicsUC,
		handler:  handler,
	}
}

func TestTopicsHandler_CreateTopic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	accessor := newTopicsHandlerAccessor(ctrl)
	topicsUC := accessor.topicsUC
	h := accessor.handler

	body := `{
			"name": "Politics",
			"description": "A topic of politics in the world",
			"slug": "politics"
		}`
	tests := []struct {
		testname  string
		body      string
		initMock  func()
		assertion func(echo.Context, *httptest.ResponseRecorder)
	}{
		{
			testname: "successful topic creation",
			body:     body,
			initMock: func() {
				topicsUC.EXPECT().CreateTopic(gomock.Any(), gomock.Any()).Return(nil)
			},
			assertion: func(ctx echo.Context, rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, rr.Code)
			},
		},
		{
			testname: "bind error - invalid json",
			body:     `{}`,
			initMock: func() {},
			assertion: func(ctx echo.Context, rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rr.Code)
			},
		},
		{
			testname: "validation error - missing required fields",
			body: `{
					"description": "A topic of politics in the world",
				}`,
			initMock: func() {},
			assertion: func(ctx echo.Context, rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rr.Code)
			},
		},
		{
			testname: "usecase error - duplicate slug",
			body:     body,
			initMock: func() {
				topicsUC.EXPECT().CreateTopic(gomock.Any(), gomock.Any()).Return(errors.New("slug already exists"))
			},
			assertion: func(ctx echo.Context, rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			tt.initMock()

			req := httptest.NewRequest(http.MethodPost, "/api/v1/topics", strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			h.CreateTopic(c)

			tt.assertion(c, rec)
		})
	}
}

func TestTopicsHandler_GetTopics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	accessor := newTopicsHandlerAccessor(ctrl)
	topicsUC := accessor.topicsUC
	h := accessor.handler

	mockTime := time.Now()
	tests := []struct {
		testname  string
		initMock  func()
		response  string
		assertion func(*httptest.ResponseRecorder, string)
	}{
		{
			testname: "successful retrieval of topics",
			initMock: func() {
				mockTopics := []response.Topic{
					{ID: 1, Name: "Topic One", Description: utils.StringPtr("Desc 1"), Slug: "topic-one", UpdatedAt: mockTime},
					{ID: 2, Name: "Topic Two", Description: nil, Slug: "topic-two", UpdatedAt: mockTime},
				}
				topicsUC.EXPECT().GetTopics(gomock.Any()).Return(mockTopics, nil)
			},
			response: `{"data":[{"id":1,"name":"Topic One","description":"Desc 1","slug":"topic-one","updated_at":"` + mockTime.Format("2006-01-02T15:04:05.999999Z07:00") + `"},{"id":2,"name":"Topic Two","description":null,"slug":"topic-two","updated_at":"` + mockTime.Format("2006-01-02T15:04:05.999999Z07:00") + `"}],"http_status":200}`,
			assertion: func(rr *httptest.ResponseRecorder, s string) {
				assert.Equal(t, http.StatusOK, rr.Code)
				assert.Equal(t, s, strings.TrimSpace(rr.Body.String()))
			},
		},
		{
			testname: "no topics found (empty slice)",
			initMock: func() {
				topicsUC.EXPECT().GetTopics(gomock.Any()).Return([]response.Topic{}, nil)
			},
			response: `{"data":[],"http_status":200}`,
			assertion: func(rr *httptest.ResponseRecorder, s string) {
				assert.Equal(t, http.StatusOK, rr.Code)
				assert.Equal(t, s, strings.TrimSpace(rr.Body.String()))
			},
		},
		{
			testname: "usecase error when getting topics",
			initMock: func() {
				topicsUC.EXPECT().GetTopics(gomock.Any()).Return(nil, errors.New("database error"))
			},
			response: `{"code":422,"status":"unprocessable entity","message":"failed get topics"}`,
			assertion: func(rr *httptest.ResponseRecorder, s string) {
				assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			tt.initMock()

			req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			h.GetTopics(c)

			tt.assertion(rec, tt.response)
		})
	}
}

func TestTopicsHandler_UpdateTopic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	accessor := newTopicsHandlerAccessor(ctrl)
	topicsUC := accessor.topicsUC
	h := accessor.handler

	tests := []struct {
		testname    string
		paramID     string
		requestBody request.UpdateTopicRequest
		initMock    func()
		response    string
		assertion   func(*httptest.ResponseRecorder, string)
	}{
		{
			testname: "successful topic update",
			paramID:  "123",
			requestBody: request.UpdateTopicRequest{
				Name: utils.StringPtr("Updated Name"),
			},
			initMock: func() {
				topicsUC.EXPECT().UpdateTopic(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			response: `{"message":"topic updated","http_status":200}`,
			assertion: func(rr *httptest.ResponseRecorder, s string) {
				assert.Equal(t, http.StatusOK, rr.Code)
				assert.Equal(t, s, strings.TrimSpace(rr.Body.String()))
			},
		},
		{
			testname: "invalid id parameter",
			paramID:  "abc", // Invalid ID
			requestBody: request.UpdateTopicRequest{
				Name: utils.StringPtr("Any Name"),
			},
			initMock: func() {},
			response: `{"message":"invalid id","http_status":400}`,
			assertion: func(rr *httptest.ResponseRecorder, s string) {
				assert.Equal(t, http.StatusBadRequest, rr.Code)
				assert.Equal(t, s, strings.TrimSpace(rr.Body.String()))
			},
		},
		{
			testname: "validation error - no fields provided for update",
			paramID:  "123",
			requestBody: request.UpdateTopicRequest{
				Name: utils.StringPtr("N"),
			},
			initMock: func() {},
			response: `{"message":"update topic require one of [name, description, slug]","http_status":400}`,
			assertion: func(rr *httptest.ResponseRecorder, s string) {
				assert.Equal(t, http.StatusBadRequest, rr.Code)
				assert.Equal(t, s, strings.TrimSpace(rr.Body.String()))
			},
		},
		{
			testname: "usecase returns ErrNoFieldUpdate (code 20005)",
			paramID:  "123",
			requestBody: request.UpdateTopicRequest{
				Name: utils.StringPtr("Same Name"), // Assume this leads to ErrNoFieldUpdate
			},
			initMock: func() {
				// Create a CustomError with Code 20005
				customErr := exception.CustomError{
					Code:    20005,
					Message: "no field updated",
				}
				topicsUC.EXPECT().UpdateTopic(gomock.Any(), gomock.Any(), gomock.Any()).Return(customErr)
			},
			response: `{"message":"no field updated","http_status":200}`,
			assertion: func(rr *httptest.ResponseRecorder, s string) {
				assert.Equal(t, http.StatusOK, rr.Code)
				assert.Equal(t, s, strings.TrimSpace(rr.Body.String()))
			},
		},
		{
			testname: "usecase returns generic error",
			paramID:  "123",
			requestBody: request.UpdateTopicRequest{
				Name: utils.StringPtr("New Name"),
			},
			initMock: func() {
				topicsUC.EXPECT().UpdateTopic(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("internal server error"))
			},
			response: `{"message":"internal server error","http_status":422}`,
			assertion: func(rr *httptest.ResponseRecorder, s string) {
				assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
				assert.Equal(t, s, strings.TrimSpace(rr.Body.String()))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			tt.initMock()

			jsonBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPatch, "/api/v1/topics/"+tt.paramID, bytes.NewReader(jsonBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.paramID)
			h.UpdateTopic(c)

			tt.assertion(rec, tt.response)
		})
	}
}
