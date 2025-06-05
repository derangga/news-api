package usecase

import (
	"context"
	"errors"
	"newsapi/internal/exception"
	"newsapi/internal/model/entity"
	"newsapi/internal/model/request"
	"newsapi/internal/model/response"
	"newsapi/internal/repository"
	"newsapi/internal/utils"

	"github.com/labstack/gommon/log"
)

type topicsUsecase struct {
	repo repository.TopicsRepository
}

func NewTopicsUsecase(
	repo repository.TopicsRepository,
) TopicsUsecase {
	return topicsUsecase{repo: repo}
}

func (u topicsUsecase) CreateTopic(
	ctx context.Context,
	body request.CreateTopicRequest,
) error {
	entity := &entity.Topic{
		Name:        body.Name,
		Description: body.Description,
		Slug:        body.Slug,
	}
	err := u.repo.Create(ctx, entity)
	if err != nil {
		if valid, hint := utils.IsDuplicateKey(err); valid {
			return errors.New(hint)
		}
		return exception.ErrFailedInsertTopic
	}

	return nil
}

func (u topicsUsecase) GetTopics(ctx context.Context) ([]response.Topic, error) {
	topics, err := u.repo.GetAll(ctx)
	if err != nil {
		log.Errorf("failed get topic: %w", err)
		return nil, exception.ErrFailedGetTopic
	}

	res := []response.Topic{}
	for _, t := range topics {
		res = append(res, response.TopicSeriliazer(t))
	}

	return res, nil
}

func (u topicsUsecase) UpdateTopic(ctx context.Context, body request.UpdateTopicRequest) error {
	currentTopic, err := u.repo.GetByID(ctx, body.ID)
	if err != nil {
		if notFound := utils.IsNoRowError(err); notFound {
			return exception.ErrTopicNotFound
		}

		log.Errorf("failed get topic: %s", err.Error())
		return exception.ErrFailedUpdateTopic
	}

	updateFields := make([]string, 0)
	updatedTopic := currentTopic

	if body.Name != nil && *body.Name != currentTopic.Name {
		updatedTopic.Name = *body.Name
		updateFields = append(updateFields, "name")
	}

	if body.Description != nil {
		if currentTopic.Description == nil || *body.Description != *currentTopic.Description {
			updatedTopic.Description = body.Description
			updateFields = append(updateFields, "description")
		}
	}

	if body.Slug != nil && *body.Slug != currentTopic.Slug {
		updatedTopic.Slug = *body.Slug
		updateFields = append(updateFields, "slug")
	}

	if len(updateFields) == 0 {
		return exception.ErrNoFieldUpdate
	}

	err = u.repo.UpdateTopicFileds(ctx, &updatedTopic, updateFields)
	if err != nil {
		if valid, hint := utils.IsDuplicateKey(err); valid {
			return errors.New(hint)
		}
		return exception.ErrFailedUpdateTopic
	}

	return nil
}
