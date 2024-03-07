package daos

import (
	"api-server/daos/models"
	"api-server/db"
	"context"
	"errors"

	"gorm.io/gorm"
)

type IDocument interface {
	Save(ctx *context.Context, document *models.Document) error
}

type Document struct {
	DB *gorm.DB
}

func NewDocument() IDocument {
	return &Document{
		DB: db.New().DB,
	}
}

func (t *Document) getTable() string {
	return "documents"
}

func (t *Document) Save(ctx *context.Context, document *models.Document) error {
	tx := t.DB.Table(t.getTable()).Where("name = ? AND url = ?", document.Name, document.URL)

	var isExists bool
	err := t.DB.Raw("SELECT EXISTS (?)", tx).Scan(&isExists).Error
	if err != nil {
		return err
	}

	if isExists {
		return errors.New("document with same name already exists. Please remove old one first")
	}

	return t.DB.Table(t.getTable()).Save(document).Error
}
