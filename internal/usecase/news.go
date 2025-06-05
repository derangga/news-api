package usecase

import (
	"context"
	"database/sql"
	"errors"
	"newsapi/internal/exception"
	"newsapi/internal/model/entity"
	"newsapi/internal/model/request"
	"newsapi/internal/model/response"
	"newsapi/internal/repository"
	"newsapi/internal/utils"
	"slices"
	"time"

	"github.com/labstack/gommon/log"
)

type newsArticlesUsecase struct {
	newsArticlesrepo repository.NewsArticlesRepository
	newsTopicsRepo   repository.NewsTopicsRepository
}

func NewNewsArticlesUsecase(
	newsArticlesrepo repository.NewsArticlesRepository,
	newsTopicsRepo repository.NewsTopicsRepository,
) NewsUsecase {
	return newsArticlesUsecase{
		newsArticlesrepo: newsArticlesrepo,
		newsTopicsRepo:   newsTopicsRepo,
	}
}

func (u newsArticlesUsecase) CreateNewsArticle(
	ctx context.Context,
	body request.CreateNewsArticleRequest,
) error {

	status := "draft"
	if body.Status != nil {
		status = *body.Status
	}

	article := &entity.NewsArticle{
		Title:    body.Title,
		Content:  body.Content,
		Summary:  body.Summary,
		AuthorID: body.AuthorID,
		Slug:     body.Slug,
		Status:   entity.ArticleStatus(status),
	}

	if status == "published" {
		article.PublishedAt = sql.NullTime{Time: time.Now(), Valid: true}
	}

	articleID, err := u.newsArticlesrepo.Create(ctx, article)
	if err != nil {
		if valid, hint := utils.IsDuplicateKey(err); valid {
			return errors.New(hint)
		}
		log.Errorf("failed create news: %w", err)
		return exception.ErrFailedInsertNews
	}
	article.ID = articleID

	if len(body.TopicIDs) > 0 {
		err = u.newsTopicsRepo.Create(ctx, articleID, body.TopicIDs)
		if err != nil {
			log.Errorf("failed create news topic: %w", err)
			return err
		}
	}

	return nil
}

func (u newsArticlesUsecase) GetNewsArticles(ctx context.Context) ([]response.NewsArticle, error) {
	newsArticles, err := u.newsArticlesrepo.GetAll(ctx)
	if err != nil {
		log.Errorf("failed get topic: %w", err)
		return nil, exception.ErrFailedGetNews
	}

	res := []response.NewsArticle{}
	for _, na := range newsArticles {
		res = append(res, response.NewsArticleSeriliazer(na))
	}

	return res, nil
}

func (u newsArticlesUsecase) GetNewsArticleBySlug(ctx context.Context, slug string) (response.NewsArticleWithTopic, error) {
	newsArticles, err := u.newsArticlesrepo.GetActiveArticleBySlug(ctx, slug)
	if err != nil {
		if notFound := utils.IsNoRowError(err); notFound {
			return response.NewsArticleWithTopic{}, exception.ErrNewsNotFound
		}

		log.Errorf("failed get news: %s", err.Error())
		return response.NewsArticleWithTopic{}, exception.ErrFailedGetNews
	}

	res := response.NewsArticleWithTopicSerializer(newsArticles)
	return res, nil
}

func (u newsArticlesUsecase) UpdateNewsArticleBySlug(
	ctx context.Context,
	slug string,
	body request.UpdateNewsArticleRequest,
) error {
	currentNews, err := u.newsArticlesrepo.GetArticleBySlug(ctx, slug)
	if err != nil {
		if notFound := utils.IsNoRowError(err); notFound {
			return exception.ErrNewsNotFound
		}

		log.Errorf("failed get news: %s", err.Error())
		return exception.ErrFailedGetNews
	}

	// Prepare update fields
	updateFields := make([]string, 0)
	updatedNews := currentNews

	if body.Title != nil && *body.Title != currentNews.Title {
		updatedNews.Title = *body.Title
		updateFields = append(updateFields, "title")
	}

	if body.Content != nil && *body.Content != currentNews.Content {
		updatedNews.Content = *body.Content
		updateFields = append(updateFields, "content")
	}

	if body.Summary != nil {
		updatedNews.Summary = body.Summary
		updateFields = append(updateFields, "summary")
	}
	if body.Slug != nil && *body.Slug != currentNews.Slug {
		updatedNews.Slug = *body.Slug
		updateFields = append(updateFields, "slug")
	}
	if body.Status != nil && *body.Status != string(currentNews.Status) {
		updatedNews.Status = entity.ArticleStatus(*body.Status)
		updateFields = append(updateFields, "status")
	}

	currentTopics := append([]int32(nil), currentNews.Topics...)
	isTopicSame := slices.Equal(body.TopicIDs, currentTopics)

	if len(updateFields) == 0 && isTopicSame {
		return exception.ErrNoFieldUpdate
	}

	err = u.newsArticlesrepo.UpdateArticleFields(ctx, &currentNews, updateFields)
	if err != nil {
		// Handle unique constraint violation
		if valid, hint := utils.IsDuplicateKey(err); valid {
			return errors.New(hint)
		}
		log.Errorf("failed update news: %s", err.Error())
		return exception.ErrFailedUpdateNews
	}

	if body.TopicIDs != nil {
		err = u.newsTopicsRepo.ReplaceArticleTopics(ctx, currentNews.ID, body.TopicIDs)
		if err != nil {
			log.Errorf("failed get news: %s", err.Error())
			return exception.ErrFailedUpdateTopicNews
		}
	}

	return nil
}

func (u newsArticlesUsecase) DeleteNewsArticleBySlug(ctx context.Context, slug string) error {
	currentNews, err := u.newsArticlesrepo.GetArticleBySlug(ctx, slug)
	if err != nil {
		if notFound := utils.IsNoRowError(err); notFound {
			return exception.ErrNewsNotFound
		}

		log.Errorf("failed get news: %s", err.Error())
		return exception.ErrFailedGetNews
	}

	err = u.newsArticlesrepo.DeleteBySlug(ctx, slug)
	if err != nil {
		log.Errorf("failed delete news: %s", err.Error())
		return exception.ErrFailedDeleteNews
	}

	err = u.newsTopicsRepo.DeleteByArticleID(ctx, currentNews.ID)
	if err != nil {
		log.Errorf("failed delete topic news: %s", err.Error())
		return exception.ErrFailedDeleteTopicNews
	}

	return nil
}
