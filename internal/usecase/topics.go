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
		res = append(res, response.Topic{
			ID:          t.ID,
			Name:        t.Name,
			Description: t.Description,
			Slug:        t.Slug,
			UpdatedAt:   t.UpdatedAt,
		})
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

	// Prepare update fields
	updateFields := make([]string, 0)
	updatedTopic := currentTopic

	// Check each field for changes
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

	// If no fields to update, return existing topic
	if len(updateFields) == 0 {
		return exception.ErrNoFieldUpdate
	}

	// Perform update
	err = u.repo.Update(ctx, &updatedTopic, updateFields)
	if err != nil {
		// Handle unique constraint violation
		if valid, hint := utils.IsDuplicateKey(err); valid {
			return errors.New(hint)
		}
		return exception.ErrFailedUpdateTopic
	}

	return nil
}
