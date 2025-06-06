package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"newsapi/internal/handler"
	mock_usecase "newsapi/mocks/usecase"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type UsersHandlerAccessor struct {
	usersUC *mock_usecase.MockUsersUsecase
	handler handler.UsersHandler
}

func newUsersHandlerAccessor(ctrl *gomock.Controller) UsersHandlerAccessor {
	usersUC := mock_usecase.NewMockUsersUsecase(ctrl)
	validator := validator.New()
	handler := handler.NewUsersHandler(validator, usersUC)
	return UsersHandlerAccessor{
		usersUC: usersUC,
		handler: handler,
	}
}

func Test_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	accessor := newUsersHandlerAccessor(ctrl)
	h := accessor.handler
	e := echo.New()
	body := `{"name":"user 1","email":"user@example.com"}`

	tests := []struct {
		name      string
		body      string
		initMock  func()
		assertion func(echo.Context, *httptest.ResponseRecorder, error)
	}{
		{
			name:     "request register and return 400",
			body:     `{"name":"user 1"}`,
			initMock: func() {},
			assertion: func(ctx echo.Context, rr *httptest.ResponseRecorder, err error) {
				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, rr.Code)
			},
		},
		{
			name: "request register and return 422",
			body: body,
			initMock: func() {
				accessor.usersUC.EXPECT().CreateUser(gomock.Any(), gomock.Any()).
					Return(errors.New("error create"))
			},
			assertion: func(ctx echo.Context, rr *httptest.ResponseRecorder, err error) {
				assert.NoError(t, err)
				assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
			},
		},
		{
			name: "request register and return 200",
			body: body,
			initMock: func() {
				accessor.usersUC.EXPECT().CreateUser(gomock.Any(), gomock.Any()).
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

			req := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := h.CreateUser(c)

			tt.assertion(c, rec, err)
		})
	}
}
