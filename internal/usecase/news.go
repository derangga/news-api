package usecase

import (
	"context"
	"newsapi/internal/model/request"
	"newsapi/internal/repository"
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

func (u newsArticlesUsecase) CreateNewsArticle(ctx context.Context, body request.CreateNewsArticleRequest) {

}
