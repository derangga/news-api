package usecase_test

import (
	"context"
	"database/sql"
	"errors"
	"newsapi/internal/exception"
	"newsapi/internal/model/entity"
	"newsapi/internal/model/request"
	"newsapi/internal/model/response"
	"newsapi/internal/usecase"
	"newsapi/internal/utils"
	mock_repository "newsapi/mocks/repository"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type TopicsAccessor struct {
	topicRepo *mock_repository.MockTopicsRepository
	topicUC   usecase.TopicsUsecase
}

func newTopicAccessor(ctrl *gomock.Controller) TopicsAccessor {
	repo := mock_repository.NewMockTopicsRepository(ctrl)
	topicUC := usecase.NewTopicsUsecase(repo)
	return TopicsAccessor{
		topicRepo: repo,
		topicUC:   topicUC,
	}
}

func Test_CreateTopic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	accessor := newTopicAccessor(ctrl)
	topicUC := accessor.topicUC
	topicRepo := accessor.topicRepo
	ctx := context.Background()

	mockReq := request.CreateTopicRequest{
		Name: "test",
		Slug: "test-1",
	}
	tests := []struct {
		testname  string
		initMock  func()
		assertion func(err error)
	}{
		{
			testname: "create topic and repository return err then uc return error",
			initMock: func() {
				topicRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("error"))
			},
			assertion: func(err error) {
				assert.Error(t, err)
			},
		},
		{
			testname: "success create topic then return err nil",
			initMock: func() {
				topicRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			assertion: func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			tt.initMock()
			err := topicUC.CreateTopic(ctx, mockReq)
			tt.assertion(err)
		})
	}
}

func Test_GetTopics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	accessor := newTopicAccessor(ctrl)
	topicUC := accessor.topicUC
	topicRepo := accessor.topicRepo
	ctx := context.Background()

	tests := []struct {
		testname  string
		initMock  func()
		assertion func(topics []response.Topic, err error)
	}{
		{
			testname: "successful retrieval of multiple topics",
			initMock: func() {
				mockTopics := []entity.Topic{
					{
						ID:          1,
						Name:        "Topic A",
						Description: utils.StringPtr("Description A"),
						Slug:        "topic-a",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
					{
						ID:          2,
						Name:        "Topic B",
						Description: utils.StringPtr("Description B"),
						Slug:        "topic-b",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
				}
				topicRepo.EXPECT().GetAll(ctx).Return(mockTopics, nil)
			},
			assertion: func(topics []response.Topic, err error) {
				assert.NoError(t, err)
				assert.Len(t, topics, 2)
				assert.Equal(t, "Topic A", topics[0].Name)
				assert.Equal(t, "topic-b", topics[1].Slug)
			},
		},
		{
			testname: "successful retrieval of no topics (empty slice)",
			initMock: func() {
				topicRepo.EXPECT().GetAll(ctx).Return([]entity.Topic{}, nil)
			},
			assertion: func(topics []response.Topic, err error) {
				assert.NoError(t, err)
				assert.Len(t, topics, 0)
				assert.Empty(t, topics)
			},
		},
		{
			testname: "repository returns an error",
			initMock: func() {
				topicRepo.EXPECT().GetAll(ctx).Return(nil, errors.New("database connection lost"))
			},
			assertion: func(topics []response.Topic, err error) {
				assert.Error(t, err)
				assert.Nil(t, topics)
				assert.Equal(t, exception.ErrFailedGetTopic, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			tt.initMock()
			topics, err := topicUC.GetTopics(ctx)
			tt.assertion(topics, err)
		})
	}
}

func Test_UpdateTopic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	accessor := newTopicAccessor(ctrl)
	topicUC := accessor.topicUC
	topicRepo := accessor.topicRepo
	ctx := context.Background()

	tests := []struct {
		testname  string
		mockID    int
		mockTopic entity.Topic
		mockReq   request.UpdateTopicRequest
		initMock  func()
		assertion func(err error)
	}{
		{
			testname: "successful update - name only",
			mockID:   1,
			mockTopic: entity.Topic{
				ID:          1,
				Name:        "Old Name",
				Description: utils.StringPtr("Old Description"),
				Slug:        "old-slug",
			},
			mockReq: request.UpdateTopicRequest{
				Name: utils.StringPtr("New Name"),
			},
			initMock: func() {
				topicRepo.EXPECT().GetByID(ctx, 1).Return(entity.Topic{
					ID:          1,
					Name:        "Old Name",
					Description: utils.StringPtr("Old Description"),
					Slug:        "old-slug",
				}, nil)
				topicRepo.EXPECT().UpdateTopicFileds(ctx, gomock.Any(), []string{"name"}).Return(nil)
			},
			assertion: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			testname: "successful update - description only (from non-nil to new non-nil)",
			mockID:   2,
			mockTopic: entity.Topic{
				ID:          2,
				Name:        "Topic 2",
				Description: utils.StringPtr("Old Desc"),
				Slug:        "topic-2",
			},
			mockReq: request.UpdateTopicRequest{
				Description: utils.StringPtr("New Description"),
			},
			initMock: func() {
				topicRepo.EXPECT().GetByID(ctx, 2).Return(entity.Topic{
					ID:          2,
					Name:        "Topic 2",
					Description: utils.StringPtr("Old Desc"),
					Slug:        "topic-2",
				}, nil)
				topicRepo.EXPECT().UpdateTopicFileds(ctx, gomock.Any(), []string{"description"}).Return(nil)
			},
			assertion: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			testname: "successful update - description only (from nil to non-nil)",
			mockID:   3,
			mockTopic: entity.Topic{
				ID:          3,
				Name:        "Topic 3",
				Description: nil,
				Slug:        "topic-3",
			},
			mockReq: request.UpdateTopicRequest{
				Description: utils.StringPtr("First Description"),
			},
			initMock: func() {
				topicRepo.EXPECT().GetByID(ctx, 3).Return(entity.Topic{
					ID:          3,
					Name:        "Topic 3",
					Description: nil,
					Slug:        "topic-3",
				}, nil)
				topicRepo.EXPECT().UpdateTopicFileds(ctx, gomock.Any(), []string{"description"}).Return(nil)
			},
			assertion: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			testname: "successful update - slug only",
			mockID:   4,
			mockTopic: entity.Topic{
				ID:          4,
				Name:        "Topic 4",
				Description: utils.StringPtr("Desc 4"),
				Slug:        "topic-4",
			},
			mockReq: request.UpdateTopicRequest{
				Slug: utils.StringPtr("new-topic-4-slug"),
			},
			initMock: func() {
				topicRepo.EXPECT().GetByID(ctx, 4).Return(entity.Topic{
					ID:          4,
					Name:        "Topic 4",
					Description: utils.StringPtr("Desc 4"),
					Slug:        "topic-4",
				}, nil)
				topicRepo.EXPECT().UpdateTopicFileds(ctx, gomock.Any(), []string{"slug"}).Return(nil)
			},
			assertion: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			testname: "successful update - multiple fields",
			mockID:   5,
			mockTopic: entity.Topic{
				ID:          5,
				Name:        "Old Name 5",
				Description: utils.StringPtr("Old Description 5"),
				Slug:        "old-slug-5",
			},
			mockReq: request.UpdateTopicRequest{
				Name:        utils.StringPtr("Updated Name 5"),
				Description: utils.StringPtr("Updated Description 5"),
				Slug:        utils.StringPtr("updated-slug-5"),
			},
			initMock: func() {
				topicRepo.EXPECT().GetByID(ctx, 5).Return(entity.Topic{
					ID:          5,
					Name:        "Old Name 5",
					Description: utils.StringPtr("Old Description 5"),
					Slug:        "old-slug-5",
				}, nil)
				topicRepo.EXPECT().UpdateTopicFileds(ctx, gomock.Any(), gomock.InAnyOrder([]string{"name", "description", "slug"})).Return(nil)
			},
			assertion: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			testname: "topic not found",
			mockID:   99,
			mockTopic: entity.Topic{
				ID: 99,
			},
			mockReq: request.UpdateTopicRequest{
				Name: utils.StringPtr("Any Name"),
			},
			initMock: func() {
				topicRepo.EXPECT().GetByID(ctx, 99).Return(entity.Topic{}, sql.ErrNoRows)
			},
			assertion: func(err error) {
				assert.Error(t, err)
				assert.Equal(t, exception.ErrTopicNotFound, err)
			},
		},
		{
			testname: "failed to get topic - generic error",
			mockID:   99,
			mockTopic: entity.Topic{
				ID: 99,
			},
			mockReq: request.UpdateTopicRequest{
				Name: utils.StringPtr("Any Name"),
			},
			initMock: func() {
				topicRepo.EXPECT().GetByID(ctx, 99).Return(entity.Topic{}, errors.New("database connection error"))
			},
			assertion: func(err error) {
				assert.Error(t, err)
				assert.Equal(t, exception.ErrFailedUpdateTopic, err)
			},
		},
		{
			testname: "no fields to update",
			mockID:   6,
			mockTopic: entity.Topic{
				ID:          6,
				Name:        "Existing Name",
				Description: utils.StringPtr("Existing Desc"),
				Slug:        "existing-slug",
			},
			mockReq: request.UpdateTopicRequest{
				Name:        utils.StringPtr("Existing Name"),
				Description: utils.StringPtr("Existing Desc"),
				Slug:        utils.StringPtr("existing-slug"),
			},
			initMock: func() {
				topicRepo.EXPECT().GetByID(ctx, 6).Return(entity.Topic{
					ID:          6,
					Name:        "Existing Name",
					Description: utils.StringPtr("Existing Desc"),
					Slug:        "existing-slug",
				}, nil)
			},
			assertion: func(err error) {
				assert.Error(t, err)
				assert.Equal(t, exception.ErrNoFieldUpdate, err)
			},
		},
		{
			testname: "duplicate slug error on update",
			mockID:   7,
			mockTopic: entity.Topic{
				ID:          7,
				Name:        "Topic 7",
				Description: utils.StringPtr("Desc 7"),
				Slug:        "topic-7",
			},
			mockReq: request.UpdateTopicRequest{
				Slug: utils.StringPtr("duplicate-slug"),
			},
			initMock: func() {
				topicRepo.EXPECT().GetByID(ctx, 7).Return(entity.Topic{
					ID:          7,
					Name:        "Topic 7",
					Description: utils.StringPtr("Desc 7"),
					Slug:        "topic-7",
				}, nil)
				topicRepo.EXPECT().UpdateTopicFileds(ctx, gomock.Any(), []string{"slug"}).
					Return(errors.New("pq: duplicate key value violates unique constraint \"topics_slug_key\""))
			},
			assertion: func(err error) {
				assert.Error(t, err)
				assert.Equal(t, "failed update topic", err.Error())
			},
		},
		{
			testname: "failed to update topic - generic error",
			mockID:   8,
			mockTopic: entity.Topic{
				ID:          8,
				Name:        "Topic 8",
				Description: utils.StringPtr("Desc 8"),
				Slug:        "topic-8",
			},
			mockReq: request.UpdateTopicRequest{
				Name: utils.StringPtr("New Name 8"),
			},
			initMock: func() {
				topicRepo.EXPECT().GetByID(ctx, 8).Return(entity.Topic{
					ID:          8,
					Name:        "Topic 8",
					Description: utils.StringPtr("Desc 8"),
					Slug:        "topic-8",
				}, nil)
				topicRepo.EXPECT().UpdateTopicFileds(ctx, gomock.Any(), []string{"name"}).Return(errors.New("db write error"))
			},
			assertion: func(err error) {
				assert.Error(t, err)
				assert.Equal(t, exception.ErrFailedUpdateTopic, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			tt.initMock()
			err := topicUC.UpdateTopic(ctx, tt.mockID, tt.mockReq)
			tt.assertion(err)
		})
	}
}

func Test_DeleteTopic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	accessor := newTopicAccessor(ctrl)
	topicUC := accessor.topicUC
	topicRepo := accessor.topicRepo
	ctx := context.Background()

	tests := []struct {
		testname  string
		id        int
		initMock  func(id int)
		assertion func(err error)
	}{
		{
			testname: "delete topic but no topic id found then return error",
			id:       1,
			initMock: func(id int) {
				topicRepo.EXPECT().GetByID(gomock.Any(), id).Return(entity.Topic{}, sql.ErrNoRows)
			},
			assertion: func(err error) {
				assert.Error(t, err)
				assert.Equal(t, exception.ErrTopicNotFound, err)
			},
		},
		{
			testname: "delete topic but db throw error then return error",
			id:       1,
			initMock: func(id int) {
				topicRepo.EXPECT().GetByID(gomock.Any(), id).Return(entity.Topic{}, errors.New("error"))
			},
			assertion: func(err error) {
				assert.Error(t, err)
				assert.Equal(t, exception.ErrFailedUpdateTopic, err)
			},
		},
		{
			testname: "delete topic failed on update delete then return error",
			id:       1,
			initMock: func(id int) {
				topicRepo.EXPECT().GetByID(gomock.Any(), id).Return(entity.Topic{
					ID:   1,
					Name: "test topic",
				}, nil)
				topicRepo.EXPECT().Delete(gomock.Any(), id).Return(errors.New("failed update"))
			},
			assertion: func(err error) {
				assert.Error(t, err)
				assert.Equal(t, exception.ErrFailedDeleteTopic, err)
			},
		},
		{
			testname: "delete topic and success then return error nil",
			id:       1,
			initMock: func(id int) {
				topicRepo.EXPECT().GetByID(gomock.Any(), id).Return(entity.Topic{
					ID:   1,
					Name: "test topic",
				}, nil)
				topicRepo.EXPECT().Delete(gomock.Any(), id).Return(nil)
			},
			assertion: func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			tt.initMock(tt.id)
			err := topicUC.DeleteTopic(ctx, tt.id)
			tt.assertion(err)
		})
	}
}
