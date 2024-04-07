package entity

import "github.com/google/uuid"

type Role struct {
	ID   uuid.UUID `json:"id" form:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name string    `json:"name" form:"name"`
}
