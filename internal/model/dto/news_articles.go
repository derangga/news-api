package dto

import "newsapi/internal/model/entity"

type NewsFilter struct {
	Status  entity.ArticleStatus
	TopicID string
}
