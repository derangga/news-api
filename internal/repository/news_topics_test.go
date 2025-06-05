package repository_test

import (
	"context"
	"errors"
	"fmt"
	"newsapi/internal/repository"
	"newsapi/internal/utils"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Test_CreateNewsTopics(t *testing.T) {
	mockDb, sqlxDB, mockSql := utils.GenerateMockDb()
	defer mockDb.Close()

	repos := repository.NewNewsTopicsRepository(sqlxDB)
	ctx := context.Background()

	t.Run(
		"returns early if no topicIDs",
		func(t *testing.T) {
			err := repos.Create(ctx, 1, []int{})
			assert.NoError(t, err)
		},
	)

	t.Run(
		"successful insert of multiple topics",
		func(t *testing.T) {
			topicIDs := []int{1, 2}
			articleID := 10

			mockSql.ExpectBegin()
			mockSql.ExpectExec(`INSERT INTO news_topics \(news_article_id, topic_id\) VALUES \(\$1, \$2\),\(\$3, \$4\)`).
				WithArgs(articleID, 1, articleID, 2).
				WillReturnResult(sqlmock.NewResult(0, 2))
			mockSql.ExpectCommit()

			err := repos.Create(ctx, articleID, topicIDs)
			assert.NoError(t, err)
		},
	)

	t.Run("foreign key constraint error", func(t *testing.T) {
		topicIDs := []int{999}
		articleID := 10

		mockSql.ExpectBegin()
		mockSql.ExpectExec(`INSERT INTO news_topics`).WithArgs(articleID, 999).
			WillReturnError(fmt.Errorf("pq: insert or update on table violates foreign key constraint"))
		mockSql.ExpectRollback()

		err := repos.Create(ctx, articleID, topicIDs)
		assert.EqualError(t, err, "one or more topic IDs are invalid")
	})
}

func Test_ReplaceArticleTopics(t *testing.T) {
	mockDb, sqlxDB, mockSql := utils.GenerateMockDb()
	defer mockDb.Close()

	repos := repository.NewNewsTopicsRepository(sqlxDB)
	ctx := context.Background()

	t.Run("successfully replaces topics", func(t *testing.T) {
		articleID := 1
		topicIDs := []int32{10, 20}

		mockSql.ExpectBegin()

		// existing topics in DB
		mockSql.ExpectQuery(`SELECT topic_id FROM news_topics`).
			WithArgs(articleID).
			WillReturnRows(sqlmock.NewRows([]string{"topic_id"}).AddRow(10).AddRow(30))

		// delete topic 30
		mockSql.ExpectExec(`UPDATE news_topics SET deleted_at = NOW\(\) WHERE news_article_id = \$1 AND topic_id = \$2 AND deleted_at IS NULL`).
			WithArgs(articleID, 30).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// upsert topic 10
		mockSql.ExpectExec(`INSERT INTO news_topics .* ON CONFLICT .* DO UPDATE SET deleted_at = NULL`).
			WithArgs(articleID, 10).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// insert new topic 20
		mockSql.ExpectExec(`INSERT INTO news_topics .* ON CONFLICT .* DO UPDATE SET deleted_at = NULL`).
			WithArgs(articleID, 20).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mockSql.ExpectCommit()

		err := repos.ReplaceArticleTopics(ctx, articleID, topicIDs)
		assert.NoError(t, err)
	})

	t.Run("select context fails", func(t *testing.T) {
		articleID := 2
		topicIDs := []int32{1, 2}

		mockSql.ExpectBegin()
		mockSql.ExpectQuery(`SELECT topic_id FROM news_topics`).
			WithArgs(articleID).
			WillReturnError(errors.New("db error"))
		mockSql.ExpectRollback()

		err := repos.ReplaceArticleTopics(ctx, articleID, topicIDs)
		assert.Error(t, err)
	})
}
