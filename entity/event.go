package entity

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID          uuid.UUID `json:"id" form:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()" `
	Name        string    `json:"name" form:"name"`
	Description string    `json:"description" form:"description" gorm:"type:text"`
	Price       int       `json:"price" form:"price"`

	Capacity  int `json:"capacity,omitempty" form:"capacity"`
	Registers int `json:"registers,omitempty" form:"registers"`

	StartDate time.Time `json:"start_date" form:"start_date" gorm:"type:timestamp without time zone;default:null"`
	EndDate   time.Time `json:"end_date" form:"end_date" gorm:"type:timestamp without time zone;default:null"`

	Participants []User `json:"participants,omitempty" gorm:"many2many:tickets;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Timestamp
}
