package usecase_test

import (
	"context"
	"errors"
	"newsapi/internal/model/request"
	"newsapi/internal/usecase"
	mock_repository "newsapi/mocks/repository"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type UsersAccessor struct {
	userRepo *mock_repository.MockUsersRepository
	userUC   usecase.UsersUsecase
}

func newUserAccessor(ctrl *gomock.Controller) UsersAccessor {
	repo := mock_repository.NewMockUsersRepository(ctrl)
	userUC := usecase.NewUsersUsecase(repo)
	return UsersAccessor{
		userRepo: repo,
		userUC:   userUC,
	}
}

func Test_CreateUse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	accessor := newUserAccessor(ctrl)
	repo := accessor.userRepo
	uc := accessor.userUC
	ctx := context.Background()

	mockReq := request.CreateUserRequest{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	tests := []struct {
		testname  string
		initMock  func()
		assertion func(err error)
	}{
		{
			testname: "create user and repository return err then uc return error",
			initMock: func() {
				repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("error"))
			},
			assertion: func(err error) {
				assert.Error(t, err)
			},
		},
		{
			testname: "success create user then return err nil",
			initMock: func() {
				repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			assertion: func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			tt.initMock()
			err := uc.CreateUser(ctx, mockReq)
			tt.assertion(err)
		})
	}
}
