package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"newsapi/internal/model/entity"
	"newsapi/internal/repository"
	"newsapi/internal/utils"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func Test_CreateNewsArticle(t *testing.T) {
	mockDb, sqlxDB, mockSql := utils.GenerateMockDb()
	defer mockDb.Close()

	repos := repository.NewNewsArticlesRepository(sqlxDB)
	ctx := context.Background()

	query := `INSERT INTO news_articles \(title, content, summary, author_id, slug, status, published_at\) VALUES \(\?, \?, \?, \?, \?, \?, \?\) RETURNING id`
	tests := []struct {
		testname  string
		entity    entity.NewsArticle
		initMock  func(entity.NewsArticle)
		assertion func(err error)
	}{
		{
			testname: "insert article then return error",
			entity: entity.NewsArticle{
				Title:    "Title",
				Content:  "Content",
				Summary:  utils.StringPtr("Summary"),
				AuthorID: 1,
				Slug:     "sample-slug",
				Status:   "published",
				PublishedAt: sql.NullTime{
					Time:  time.Now(),
					Valid: true,
				},
			},
			initMock: func(na entity.NewsArticle) {
				mockSql.ExpectPrepare(query).ExpectQuery().
					WithArgs(na.Title, na.Content, na.Summary, na.AuthorID, na.Slug, na.Status, na.PublishedAt).
					WillReturnError(errors.New("failed insert"))
			},
			assertion: func(err error) {
				assert.Error(t, err)
			},
		},
		{
			testname: "insert article successfully",
			entity: entity.NewsArticle{
				Title:    "Title",
				Content:  "Content",
				Summary:  utils.StringPtr("Summary"),
				AuthorID: 1,
				Slug:     "sample-slug",
				Status:   "published",
				PublishedAt: sql.NullTime{
					Time:  time.Now(),
					Valid: true,
				},
			},
			initMock: func(na entity.NewsArticle) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(10)
				mockSql.ExpectPrepare(query).ExpectQuery().
					WithArgs(na.Title, na.Content, na.Summary, na.AuthorID, na.Slug, na.Status, na.PublishedAt).
					WillReturnRows(rows)
			},
			assertion: func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			tt.initMock(tt.entity)
			_, err := repos.Create(ctx, &tt.entity)
			tt.assertion(err)
		})
	}
}

func Test_GetArticleBySlug(t *testing.T) {
	mockDb, sqlxDB, mockSql := utils.GenerateMockDb()
	defer mockDb.Close()

	repos := repository.NewNewsArticlesRepository(sqlxDB)
	ctx := context.Background()

	query := `
			SELECT a.id, a.title, a.content, a.summary, a.author_id, a.slug, 
			a.status, a.published_at, a.created_at, a.updated_at,
			COALESCE\(array_agg\(nt.topic_id ORDER BY nt.topic_id\), '\{\}'\) AS topic_ids 
			FROM news_articles a
			INNER JOIN news_topics nt ON nt.news_article_id = a.id AND nt.deleted_at IS NULL
			WHERE slug = \$1 AND a.deleted_at is NULL
			GROUP BY a.id`

	tests := []struct {
		testname  string
		slug      string
		initMock  func(slug string)
		assertion func(entity entity.NewsArticleWithTopic, err error)
	}{
		{
			testname: "get article successfully",
			slug:     "slug-123",
			initMock: func(slug string) {
				rows := sqlmock.NewRows([]string{
					"id", "title", "content", "summary", "author_id", "slug",
					"status", "published_at", "created_at", "updated_at", "topic_ids",
				}).AddRow(1, "Title", "Content", "Summary", 10, slug, "published", time.Now(), time.Now(), time.Now(), pq.Int64Array{1, 2})

				mockSql.ExpectQuery(query).WithArgs(slug).WillReturnRows(rows)
			},
			assertion: func(result entity.NewsArticleWithTopic, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "slug-123", result.Slug)
			},
		},
		{
			testname: "get article returns error",
			slug:     "missing-slug",
			initMock: func(slug string) {
				mockSql.ExpectQuery(query).WithArgs(slug).WillReturnError(sql.ErrNoRows)
			},
			assertion: func(result entity.NewsArticleWithTopic, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			tt.initMock(tt.slug)
			result, err := repos.GetArticleBySlug(ctx, tt.slug)
			tt.assertion(result, err)
		})
	}
}

func Test_GetActiveArticleBySlug(t *testing.T) {
	mockDb, sqlxDB, mockSql := utils.GenerateMockDb()
	defer mockDb.Close()

	repos := repository.NewNewsArticlesRepository(sqlxDB)
	ctx := context.Background()

	query := `SELECT
				a.id,
				a.title,
				a.content,
				a.slug,
				a.published_at,
				u.name,
				COALESCE\(array_agg\(t.name ORDER BY t.name\) FILTER \(WHERE t.name IS NOT NULL\), '\{\}'\) AS topics
			FROM news_articles a
			INNER JOIN users u on u.id = a.author_id 
			INNER JOIN news_topics nt ON nt.news_article_id = a.id AND nt.deleted_at IS NULL
			INNER JOIN topics t ON t.id = nt.topic_id
			WHERE a.slug = \$1 AND a.deleted_at IS NULL
			GROUP BY a.id, u.name`

	tests := []struct {
		testname  string
		slug      string
		initMock  func(slug string)
		assertion func(entity entity.ActiveNewsWithTopic, err error)
	}{
		{
			testname: "get active article successfully",
			slug:     "active-slug",
			initMock: func(slug string) {
				rows := sqlmock.NewRows([]string{
					"id", "title", "content", "slug", "published_at", "name", "topics",
				}).AddRow(1, "Active Title", "Content", slug, time.Now(), "Author Name", pq.StringArray{"Tech", "Go"})

				mockSql.ExpectQuery(query).WithArgs(slug).WillReturnRows(rows)
			},
			assertion: func(result entity.ActiveNewsWithTopic, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "active-slug", result.Slug)
			},
		},
		{
			testname: "get active article returns error",
			slug:     "missing-article",
			initMock: func(slug string) {
				mockSql.ExpectQuery(query).WithArgs(slug).WillReturnError(sql.ErrNoRows)
			},
			assertion: func(result entity.ActiveNewsWithTopic, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			tt.initMock(tt.slug)
			result, err := repos.GetActiveArticleBySlug(ctx, tt.slug)
			tt.assertion(result, err)
		})
	}
}

func Test_UpdateArticleFields(t *testing.T) {
	mockDb, sqlxDB, mockSql := utils.GenerateMockDb()
	defer mockDb.Close()

	repos := repository.NewNewsArticlesRepository(sqlxDB)
	ctx := context.Background()

	tests := []struct {
		testname     string
		news         entity.NewsArticleWithTopic
		updateFields []string
		initMock     func(news entity.NewsArticleWithTopic, updateFields []string)
		assertion    func(err error)
	}{
		{
			testname: "update title and content successfully",
			news: entity.NewsArticleWithTopic{
				ID:      1,
				Title:   "New Title",
				Content: "Updated content",
			},
			updateFields: []string{"title", "content"},
			initMock: func(news entity.NewsArticleWithTopic, updateFields []string) {
				// Generate expected query and args
				query := `UPDATE news_articles SET title = \$1, content = \$2, published_at = \$3, updated_at = \$4 WHERE id = \$5`
				mockSql.ExpectExec(query).
					WithArgs(
						news.Title,
						news.Content,
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						news.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			assertion: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			testname: "update slug returns error",
			news: entity.NewsArticleWithTopic{
				ID:   2,
				Slug: "new-slug",
			},
			updateFields: []string{"slug"},
			initMock: func(news entity.NewsArticleWithTopic, updateFields []string) {
				query := `UPDATE news_articles SET slug = \$1, published_at = \$2, updated_at = \$3 WHERE id = \$4`
				mockSql.ExpectExec(query).
					WithArgs(
						news.Slug,
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						news.ID,
					).
					WillReturnError(errors.New("db error"))
			},
			assertion: func(err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			tt.initMock(tt.news, tt.updateFields)
			err := repos.UpdateArticleFields(ctx, &tt.news, tt.updateFields)
			tt.assertion(err)
		})
	}
}

func Test_DeleteBySlug(t *testing.T) {
	mockDb, sqlxDB, mockSql := utils.GenerateMockDb()
	defer mockDb.Close()

	repos := repository.NewNewsArticlesRepository(sqlxDB)
	ctx := context.Background()

	query := `UPDATE news_articles SET deleted_at = NOW\(\) WHERE slug = \$1 AND deleted_at IS NULL`

	tests := []struct {
		testname  string
		slug      string
		initMock  func(slug string)
		assertion func(err error)
	}{
		{
			testname: "delete article successfully",
			slug:     "to-be-deleted",
			initMock: func(slug string) {
				mockSql.ExpectExec(query).
					WithArgs(slug).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			assertion: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			testname: "delete article fails",
			slug:     "nonexistent",
			initMock: func(slug string) {
				mockSql.ExpectExec(query).
					WithArgs(slug).
					WillReturnError(errors.New("db error"))
			},
			assertion: func(err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			tt.initMock(tt.slug)
			err := repos.DeleteBySlug(ctx, tt.slug)
			tt.assertion(err)
		})
	}
}
