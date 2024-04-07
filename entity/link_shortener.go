package entity

import "github.com/google/uuid"

type LinkShortener struct {
	ID    uuid.UUID `json:"id" form:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()" `
	Alias string    `json:"alias" form:"alias"`
	Link  string    `json:"link" form:"link"`

	Timestamp
}
