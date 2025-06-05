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
	mock_repository "newsapi/mocks/repository"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type NewsAccessor struct {
	newsArticleRepo *mock_repository.MockNewsArticlesRepository
	newsTopicsRepo  *mock_repository.MockNewsTopicsRepository
	uc              usecase.NewsUsecase
}

func newNewsAccessor(ctrl *gomock.Controller) NewsAccessor {
	newsArticleRepo := mock_repository.NewMockNewsArticlesRepository(ctrl)
	newsTopicsRepo := mock_repository.NewMockNewsTopicsRepository(ctrl)
	uc := usecase.NewNewsArticlesUsecase(newsArticleRepo, newsTopicsRepo)
	return NewsAccessor{
		newsArticleRepo: newsArticleRepo,
		newsTopicsRepo:  newsTopicsRepo,
		uc:              uc,
	}
}

func Test_CreateNewsArticle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	accessor := newNewsAccessor(ctrl)
	newsArticleRepo := accessor.newsArticleRepo
	newsTopicsRepo := accessor.newsTopicsRepo
	uc := accessor.uc
	ctx := context.Background()

	// Helper function for pointer to string
	stringPtr := func(s string) *string { return &s }

	tests := []struct {
		testname  string
		mockReq   request.CreateNewsArticleRequest
		initMock  func()
		assertion func(err error)
	}{
		{
			testname: "successful creation - draft status, no topics",
			mockReq: request.CreateNewsArticleRequest{
				Title:    "Test Article 1",
				Content:  "Content of test article 1",
				Summary:  stringPtr("Summary 1"),
				AuthorID: 1,
				Slug:     "test-article-1",
				Status:   stringPtr("draft"),
				TopicIDs: []int{},
			},
			initMock: func() {
				newsArticleRepo.EXPECT().Create(ctx, gomock.Any()).Return(1, nil)
			},
			assertion: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			testname: "successful creation - published status, with topics",
			mockReq: request.CreateNewsArticleRequest{
				Title:    "Published Article",
				Content:  "Published content",
				Summary:  nil,
				AuthorID: 2,
				Slug:     "published-article",
				Status:   stringPtr("published"),
				TopicIDs: []int{10, 20, 30},
			},
			initMock: func() {
				newsArticleRepo.EXPECT().Create(ctx, gomock.Any()).Return(2, nil)
				newsTopicsRepo.EXPECT().Create(ctx, 2, []int{10, 20, 30}).Return(nil)
			},
			assertion: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			testname: "successful creation - default status (draft), with topics",
			mockReq: request.CreateNewsArticleRequest{
				Title:    "Default Status Article",
				Content:  "Content of default status",
				AuthorID: 3,
				Slug:     "default-status-article",
				Status:   nil, // Default to draft
				TopicIDs: []int{40, 50},
			},
			initMock: func() {
				newsArticleRepo.EXPECT().Create(ctx, gomock.Any()).Return(3, nil)
				newsTopicsRepo.EXPECT().Create(ctx, 3, []int{40, 50}).Return(nil)
			},
			assertion: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			testname: "failed to create news article - duplicate slug",
			mockReq: request.CreateNewsArticleRequest{
				Title:    "Duplicate Slug Article",
				Content:  "Content",
				AuthorID: 4,
				Slug:     "duplicate-slug-value",
				Status:   stringPtr("draft"),
			},
			initMock: func() {
				newsArticleRepo.EXPECT().Create(ctx, gomock.Any()).
					Return(0, errors.New(`pq: duplicate key value violates unique constraint "news_articles_slug_key"`))
			},
			assertion: func(err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed insert news")
			},
		},
		{
			testname: "failed to create news article - generic repository error",
			mockReq: request.CreateNewsArticleRequest{
				Title:    "Error Article",
				Content:  "Content",
				AuthorID: 5,
				Slug:     "error-article",
				Status:   stringPtr("draft"),
			},
			initMock: func() {
				newsArticleRepo.EXPECT().Create(ctx, gomock.Any()).
					Return(0, errors.New("database connection failed"))
			},
			assertion: func(err error) {
				assert.Error(t, err)
				assert.Equal(t, exception.ErrFailedInsertNews, err)
			},
		},
		{
			testname: "failed to create news topics after article creation",
			mockReq: request.CreateNewsArticleRequest{
				Title:    "Article with Topic Error",
				Content:  "Content",
				AuthorID: 6,
				Slug:     "article-topic-error",
				Status:   stringPtr("draft"),
				TopicIDs: []int{1, 2},
			},
			initMock: func() {
				newsArticleRepo.EXPECT().Create(ctx, gomock.Any()).Return(6, nil)
				newsTopicsRepo.EXPECT().Create(
					ctx,
					6,
					[]int{1, 2},
				).Return(errors.New("failed to insert topic relations"))
			},
			assertion: func(err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, "failed to insert topic relations")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			tt.initMock()
			err := uc.CreateNewsArticle(ctx, tt.mockReq)
			tt.assertion(err)
		})
	}
}

func Test_GetNewsArticles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	accessor := newNewsAccessor(ctrl)
	newsArticleRepo := accessor.newsArticleRepo
	uc := accessor.uc
	ctx := context.Background()

	tests := []struct {
		testname  string
		initMock  func()
		assertion func(res []response.NewsArticle, err error)
	}{
		{
			testname: "get articles repo return error then usecase return error",
			initMock: func() {
				newsArticleRepo.EXPECT().GetAll(gomock.Any()).Return(nil, errors.New("failed retrieve"))
			},
			assertion: func(res []response.NewsArticle, err error) {
				assert.Error(t, err)
				assert.Equal(t, 0, len(res))
			},
		},
		{
			testname: "get articles valid data",
			initMock: func() {
				newsArticleRepo.EXPECT().GetAll(gomock.Any()).Return(
					[]entity.NewsArticle{
						{
							ID:        1,
							Title:     "Old Title",
							Content:   "Old Content",
							Summary:   stringPtr("Old Summary"),
							Slug:      "old-slug-1",
							Status:    entity.StatusDraft,
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						},
					},
					nil,
				)
			},
			assertion: func(res []response.NewsArticle, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 1, len(res))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			tt.initMock()
			res, err := uc.GetNewsArticles(ctx)
			tt.assertion(res, err)
		})
	}
}

func Test_GetNewsArticleBySlug(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	accessor := newNewsAccessor(ctrl)
	newsArticleRepo := accessor.newsArticleRepo
	uc := accessor.uc
	ctx := context.Background()

	tests := []struct {
		testname  string
		slug      string
		initMock  func()
		assertion func(res response.NewsArticleWithTopic, err error)
	}{
		{
			testname: "successful retrieval of an active news article by slug",
			slug:     "active-article-slug",
			initMock: func() {
				mockEntity := entity.ActiveNewsWithTopic{
					ID:          1,
					Title:       "Active News Title",
					Content:     "Some active content.",
					Slug:        "active-article-slug",
					PublishedAt: sql.NullTime{Time: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC), Valid: true},
					Topics:      []string{"technology"},
				}
				newsArticleRepo.EXPECT().GetActiveArticleBySlug(ctx, "active-article-slug").Return(mockEntity, nil)
			},
			assertion: func(res response.NewsArticleWithTopic, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 1, res.ID)
				assert.Equal(t, "Active News Title", res.Title)
				assert.Equal(t, "active-article-slug", res.Slug)
				assert.Len(t, res.Topics, 1)
			},
		},
		{
			testname: "news article not found by slug",
			slug:     "non-existent-slug",
			initMock: func() {
				newsArticleRepo.EXPECT().GetActiveArticleBySlug(ctx, "non-existent-slug").Return(entity.ActiveNewsWithTopic{}, sql.ErrNoRows)
			},
			assertion: func(res response.NewsArticleWithTopic, err error) {
				assert.Error(t, err)
				assert.Equal(t, exception.ErrNewsNotFound, err)
				assert.Empty(t, res)
			},
		},
		{
			testname: "failed to get news article - generic error from repository",
			slug:     "error-getting-slug",
			initMock: func() {
				newsArticleRepo.EXPECT().GetActiveArticleBySlug(ctx, "error-getting-slug").Return(entity.ActiveNewsWithTopic{}, errors.New("unexpected database error"))
			},
			assertion: func(res response.NewsArticleWithTopic, err error) {
				assert.Error(t, err)
				assert.Equal(t, exception.ErrFailedGetNews, err)
				assert.Empty(t, res)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			tt.initMock()
			res, err := uc.GetNewsArticleBySlug(ctx, tt.slug)
			tt.assertion(res, err)
		})
	}
}

func Test_UpdateNewsArticleBySlug(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	accessor := newNewsAccessor(ctrl)
	newsArticleRepo := accessor.newsArticleRepo
	newsTopicsRepo := accessor.newsTopicsRepo
	uc := accessor.uc
	ctx := context.Background()

	// Helper function for pointer to string
	stringPtr := func(s string) *string { return &s }

	tests := []struct {
		testname  string
		slug      string
		mockReq   request.UpdateNewsArticleRequest
		initMock  func()
		assertion func(err error)
	}{
		{
			testname: "successful update - title only",
			slug:     "old-slug-1",
			mockReq: request.UpdateNewsArticleRequest{
				Title: stringPtr("New Title"),
			},
			initMock: func() {
				newsArticleRepo.EXPECT().GetArticleBySlug(ctx, "old-slug-1").Return(entity.NewsArticleWithTopic{
					ID:        1,
					Title:     "Old Title",
					Content:   "Old Content",
					Summary:   stringPtr("Old Summary"),
					Slug:      "old-slug-1",
					Status:    entity.StatusDraft,
					Topics:    []int32{1, 2},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)
				newsArticleRepo.EXPECT().UpdateArticleFields(ctx, gomock.Any(), []string{"title"}).Return(nil)
			},
			assertion: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			testname: "successful update - content and summary",
			slug:     "old-slug-2",
			mockReq: request.UpdateNewsArticleRequest{
				Content: stringPtr("New Content 2"),
				Summary: stringPtr("New Summary 2"),
			},
			initMock: func() {
				newsArticleRepo.EXPECT().GetArticleBySlug(ctx, "old-slug-2").Return(entity.NewsArticleWithTopic{
					ID:        2,
					Title:     "Title 2",
					Content:   "Content 2",
					Summary:   stringPtr("Summary 2"),
					Slug:      "old-slug-2",
					Status:    entity.StatusPublished,
					Topics:    []int32{3},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)
				newsArticleRepo.EXPECT().UpdateArticleFields(ctx, gomock.Any(), gomock.InAnyOrder([]string{"content", "summary"})).Return(nil)
			},
			assertion: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			testname: "successful update - slug and status",
			slug:     "old-slug-3",
			mockReq: request.UpdateNewsArticleRequest{
				Slug:   stringPtr("new-slug-3"),
				Status: stringPtr(string(entity.StatusPublished)),
			},
			initMock: func() {
				newsArticleRepo.EXPECT().GetArticleBySlug(ctx, "old-slug-3").Return(entity.NewsArticleWithTopic{
					ID:        3,
					Title:     "Title 3",
					Content:   "Content 3",
					Summary:   stringPtr("Summary 3"),
					Slug:      "old-slug-3",
					Status:    entity.StatusDraft,
					Topics:    []int32{},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)
				newsArticleRepo.EXPECT().UpdateArticleFields(ctx, gomock.Any(), gomock.InAnyOrder([]string{"slug", "status"})).Return(nil)
			},
			assertion: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			testname: "successful update - topics only (same fields, different topics)",
			slug:     "slug-4",
			mockReq: request.UpdateNewsArticleRequest{
				TopicIDs: []int32{30, 40, 50},
			},
			initMock: func() {
				newsArticleRepo.EXPECT().GetArticleBySlug(ctx, "slug-4").Return(entity.NewsArticleWithTopic{
					ID:        4,
					Title:     "Title 4",
					Content:   "Content 4",
					Summary:   stringPtr("Summary 4"),
					Slug:      "slug-4",
					Status:    entity.StatusDraft,
					Topics:    []int32{10, 20},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)

				newsTopicsRepo.EXPECT().ReplaceArticleTopics(ctx, 4, []int32{30, 40, 50}).Return(nil)
			},
			assertion: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			testname: "successful update - all fields and topics",
			slug:     "slug-5",
			mockReq: request.UpdateNewsArticleRequest{
				Title:    stringPtr("New Title 5"),
				Content:  stringPtr("New Content 5"),
				Summary:  stringPtr("New Summary 5"),
				Slug:     stringPtr("new-slug-5"),
				Status:   stringPtr(string(entity.StatusPublished)),
				TopicIDs: []int32{100, 200},
			},
			initMock: func() {
				newsArticleRepo.EXPECT().GetArticleBySlug(ctx, "slug-5").Return(entity.NewsArticleWithTopic{
					ID:        5,
					Title:     "Old Title 5",
					Content:   "Old Content 5",
					Summary:   stringPtr("Old Summary 5"),
					Slug:      "slug-5",
					Status:    entity.StatusDraft,
					Topics:    []int32{1, 2},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)
				newsArticleRepo.EXPECT().UpdateArticleFields(ctx, gomock.Any(), gomock.InAnyOrder([]string{"title", "content", "summary", "slug", "status"})).Return(nil)
				newsTopicsRepo.EXPECT().ReplaceArticleTopics(ctx, 5, []int32{100, 200}).Return(nil)
			},
			assertion: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			testname: "news article not found by slug",
			slug:     "non-existent-slug",
			mockReq: request.UpdateNewsArticleRequest{
				Title: stringPtr("Any Title"),
			},
			initMock: func() {
				newsArticleRepo.EXPECT().GetArticleBySlug(ctx, "non-existent-slug").Return(entity.NewsArticleWithTopic{}, sql.ErrNoRows)
			},
			assertion: func(err error) {
				assert.Error(t, err)
				assert.Equal(t, exception.ErrNewsNotFound, err)
			},
		},
		{
			testname: "failed to get news article - generic error",
			slug:     "error-slug",
			mockReq: request.UpdateNewsArticleRequest{
				Title: stringPtr("Any Title"),
			},
			initMock: func() {
				newsArticleRepo.EXPECT().GetArticleBySlug(ctx, "error-slug").Return(entity.NewsArticleWithTopic{}, errors.New("database read error"))
			},
			assertion: func(err error) {
				assert.Error(t, err)
				assert.Equal(t, exception.ErrFailedGetNews, err)
			},
		},
		{
			testname: "no fields to update and same topics",
			slug:     "slug-6",
			mockReq: request.UpdateNewsArticleRequest{
				Title:    stringPtr("Existing Title"),
				Content:  stringPtr("Existing Content"),
				Summary:  stringPtr("Existing Summary"),
				Slug:     stringPtr("slug-6"),
				Status:   stringPtr(string(entity.StatusDraft)),
				TopicIDs: []int32{1, 2},
			},
			initMock: func() {
				newsArticleRepo.EXPECT().GetArticleBySlug(ctx, "slug-6").Return(entity.NewsArticleWithTopic{
					ID:        6,
					Title:     "Existing Title",
					Content:   "Existing Content",
					Summary:   stringPtr("Existing Summary"),
					Slug:      "slug-6",
					Status:    entity.StatusDraft,
					Topics:    []int32{1, 2},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)
			},
			assertion: func(err error) {
				assert.Error(t, err)
				assert.Equal(t, exception.ErrNoFieldUpdate, err)
			},
		},
		{
			testname: "duplicate slug error on UpdateArticleFields",
			slug:     "slug-7",
			mockReq: request.UpdateNewsArticleRequest{
				Slug: stringPtr("duplicate-slug-value"),
			},
			initMock: func() {
				newsArticleRepo.EXPECT().GetArticleBySlug(ctx, "slug-7").Return(entity.NewsArticleWithTopic{
					ID:        7,
					Title:     "Title 7",
					Content:   "Content 7",
					Summary:   stringPtr("Summary 7"),
					Slug:      "slug-7",
					Status:    entity.StatusDraft,
					Topics:    []int32{},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)
				newsArticleRepo.EXPECT().UpdateArticleFields(ctx, gomock.Any(), []string{"slug"}).Return(errors.New(`pq: duplicate key value violates unique constraint "news_articles_slug_key"`))
			},
			assertion: func(err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed update news")
			},
		},
		{
			testname: "failed to update news article fields - generic error",
			slug:     "slug-8",
			mockReq: request.UpdateNewsArticleRequest{
				Title: stringPtr("New Title 8"),
			},
			initMock: func() {
				newsArticleRepo.EXPECT().GetArticleBySlug(ctx, "slug-8").Return(entity.NewsArticleWithTopic{
					ID:        8,
					Title:     "Title 8",
					Content:   "Content 8",
					Summary:   stringPtr("Summary 8"),
					Slug:      "slug-8",
					Status:    entity.StatusDraft,
					Topics:    []int32{},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)
				newsArticleRepo.EXPECT().UpdateArticleFields(ctx, gomock.Any(), []string{"title"}).Return(errors.New("unexpected database error"))
			},
			assertion: func(err error) {
				assert.Error(t, err)
				assert.Equal(t, exception.ErrFailedUpdateNews, err)
			},
		},
		{
			testname: "failed to replace article topics",
			slug:     "slug-9",
			mockReq: request.UpdateNewsArticleRequest{
				TopicIDs: []int32{3, 4}, // Changed topics
			},
			initMock: func() {
				newsArticleRepo.EXPECT().GetArticleBySlug(ctx, "slug-9").Return(entity.NewsArticleWithTopic{
					ID:        9,
					Title:     "Title 9",
					Content:   "Content 9",
					Summary:   stringPtr("Summary 9"),
					Slug:      "slug-9",
					Status:    entity.StatusDraft,
					Topics:    []int32{1, 2},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)
				// No UpdateArticleFields if only topics change and no other fields are changed
				newsTopicsRepo.EXPECT().ReplaceArticleTopics(ctx, 9, []int32{3, 4}).Return(errors.New("failed to update topic relations"))
			},
			assertion: func(err error) {
				assert.Error(t, err)
				assert.Equal(t, exception.ErrFailedUpdateTopicNews, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			tt.initMock()
			err := uc.UpdateNewsArticleBySlug(ctx, tt.slug, tt.mockReq)
			tt.assertion(err)
		})
	}
}

func strToPtr(s string) *string {
	return &s
}

func Test_DeleteNewsArticleBySlug(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	accessor := newNewsAccessor(ctrl)
	newsArticleRepo := accessor.newsArticleRepo
	newsTopicsRepo := accessor.newsTopicsRepo
	uc := accessor.uc
	ctx := context.Background()

	tests := []struct {
		testname  string
		slug      string
		mockNews  entity.NewsArticle // News article to be returned by GetArticleBySlug
		initMock  func()
		assertion func(err error)
	}{
		{
			testname: "successful deletion of news article and associated topics",
			slug:     "article-to-delete",
			mockNews: entity.NewsArticle{
				ID:        1,
				Title:     "Article to Delete",
				Slug:      "article-to-delete",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			initMock: func() {
				newsArticleRepo.EXPECT().GetArticleBySlug(ctx, "article-to-delete").Return(entity.NewsArticleWithTopic{
					ID:    1,
					Slug:  "article-to-delete",
					Title: "Article to Delete",
				}, nil)
				newsArticleRepo.EXPECT().DeleteBySlug(ctx, "article-to-delete").Return(nil)
				newsTopicsRepo.EXPECT().DeleteByArticleID(ctx, 1).Return(nil)
			},
			assertion: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			testname: "news article not found by slug",
			slug:     "non-existent-article",
			mockNews: entity.NewsArticle{},
			initMock: func() {
				newsArticleRepo.EXPECT().GetArticleBySlug(ctx, "non-existent-article").Return(entity.NewsArticleWithTopic{}, sql.ErrNoRows)
				// No further calls to DeleteBySlug or DeleteByArticleID should occur
			},
			assertion: func(err error) {
				assert.Error(t, err)
				assert.Equal(t, exception.ErrNewsNotFound, err)
			},
		},
		{
			testname: "failed to get news article - generic error",
			slug:     "get-error-article",
			mockNews: entity.NewsArticle{}, // Not used
			initMock: func() {
				// Expect GetArticleBySlug to return a generic error
				newsArticleRepo.EXPECT().GetArticleBySlug(ctx, "get-error-article").Return(entity.NewsArticleWithTopic{}, errors.New("database connection lost"))
				// No further calls should occur
			},
			assertion: func(err error) {
				assert.Error(t, err)
				assert.Equal(t, exception.ErrFailedGetNews, err)
			},
		},
		{
			testname: "failed to delete news article by slug (article repo error)",
			slug:     "delete-error-article",
			mockNews: entity.NewsArticle{
				ID:        2,
				Title:     "Article with Delete Error",
				Slug:      "delete-error-article",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			initMock: func() {
				newsArticleRepo.EXPECT().GetArticleBySlug(ctx, "delete-error-article").Return(entity.NewsArticleWithTopic{
					ID:    2,
					Slug:  "delete-error-article",
					Title: "Article with Delete Error",
				}, nil)
				newsArticleRepo.EXPECT().DeleteBySlug(ctx, "delete-error-article").Return(errors.New("failed to delete from article table"))
			},
			assertion: func(err error) {
				assert.Error(t, err)
				assert.Equal(t, exception.ErrFailedDeleteNews, err)
			},
		},
		{
			testname: "failed to delete news topics (topics repo error)",
			slug:     "topic-delete-error-article",
			mockNews: entity.NewsArticle{
				ID:        3,
				Title:     "Article with Topic Delete Error",
				Slug:      "topic-delete-error-article",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			initMock: func() {
				// Expect GetArticleBySlug to succeed
				newsArticleRepo.EXPECT().GetArticleBySlug(ctx, "topic-delete-error-article").Return(entity.NewsArticleWithTopic{
					ID:    3,
					Slug:  "topic-delete-error-article",
					Title: "Article with Topic Delete Error",
				}, nil)
				// Expect DeleteBySlug to succeed
				newsArticleRepo.EXPECT().DeleteBySlug(ctx, "topic-delete-error-article").Return(nil)
				// Expect DeleteByArticleID to return an error
				newsTopicsRepo.EXPECT().DeleteByArticleID(ctx, 3).Return(errors.New("failed to delete topic relations"))
			},
			assertion: func(err error) {
				assert.Error(t, err)
				assert.Equal(t, exception.ErrFailedDeleteTopicNews, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			tt.initMock()
			err := uc.DeleteNewsArticleBySlug(ctx, tt.slug)
			tt.assertion(err)
		})
	}
}
