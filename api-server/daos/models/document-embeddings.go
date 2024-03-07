package models

import (
	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

type DocumentEmbedding struct {
	Id         int64
	DocumentId uuid.UUID
	Content    string
	Vector     pgvector.Vector `gorm:"type:vector(3)"`
}
