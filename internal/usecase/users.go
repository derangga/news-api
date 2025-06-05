package usecase

import (
	"context"
	"newsapi/internal/exception"
	"newsapi/internal/model/entity"
	"newsapi/internal/model/request"
	"newsapi/internal/repository"

	"github.com/labstack/gommon/log"
)

type usersUsecase struct {
	repo repository.UsersRepository
}

func NewUsersUsecase(repo repository.UsersRepository) UsersUsecase {
	return usersUsecase{repo: repo}
}

func (u usersUsecase) CreateUser(ctx context.Context, body request.CreateUserRequest) error {
	entity := &entity.User{
		Name:  body.Name,
		Email: body.Email,
	}
	err := u.repo.Create(ctx, entity)
	if err != nil {
		log.Errorf("failed create user: %w", err)
		return exception.ErrFailedInsertTopic
	}

	return nil
}
