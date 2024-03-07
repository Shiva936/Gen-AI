package daos

import (
	"api-server/daos/models"
	"api-server/db"
	"context"

	"gorm.io/gorm"
)

type IDocumentEmbedding interface {
	Save(ctx *context.Context, embedding ...*models.DocumentEmbedding) error
}

type DocumentEmbedding struct {
	DB *gorm.DB
}

func NewDocumentEmbedding() IDocumentEmbedding {
	return &DocumentEmbedding{
		DB: db.New().DB,
	}
}

func (t *DocumentEmbedding) getTable() string {
	return "documents"
}

func (t *DocumentEmbedding) Save(ctx *context.Context, embedding ...*models.DocumentEmbedding) error {
	return t.DB.Table(t.getTable()).Save(embedding).Error
}

func (t *DocumentEmbedding) Get(ctx *context.Context, queryEmbedding, similarityThreshold, numMatches string) ([]string, error) {
	if similarityThreshold == "" {
		similarityThreshold = "0.8"
	}
	if numMatches == "" {
		numMatches = "3"
	}
	res := []string{}
	err := t.DB.Table(t.getTable()).Select("content").Where("1 - (embedding <=> ?) > ? LIMIT ?)", queryEmbedding, similarityThreshold, numMatches).Scan(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}
