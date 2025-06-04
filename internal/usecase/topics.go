package usecase

import (
	"context"
	"newsapi/internal/model/request"
	"newsapi/internal/repository"
)

type topicsUsecase struct {
	repo repository.TopicsRepository
}

func NewTopicsUsecase(
	repo repository.TopicsRepository,
) TopicsUsecase {
	return topicsUsecase{repo: repo}
}

func (u topicsUsecase) CreateTopic(ctx context.Context, body request.CreateTopicRequest) {

}
