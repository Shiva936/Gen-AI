package models

import "github.com/google/uuid"

type Document struct {
	Id   uuid.UUID
	Name string
	URL  string
}
