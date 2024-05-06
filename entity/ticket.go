package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Ticket struct {
	TicketID string `json:"ticket_id" form:"ticket_id" gorm:"primaryKey" `
	UserID   string `json:"user_id" form:"user_id"`
	EventID  string `json:"event_id" form:"event_id"`

	Handphone string    `json:"handphone" form:"handphone"`
	Birthdate time.Time `json:"birthdate" form:"birthdate"`
	Seat      string    `json:"seat" form:"seat"`
	Payment   string    `json:"payment" form:"payment"`

	PaymentConfirmed *bool `json:"payment_confirmed" form:"payment_confirmed" default:"false"`
	CheckedIn        *bool `json:"checked_in" form:"checked_in" default:"false"`

	User  *User  `gorm:"foreignKey:UserID"`
	Event *Event `gorm:"foreignKey:EventID"`

	Timestamp
}

func (t *Ticket) BeforeCreate(tx *gorm.DB) error {
	var event Event
	if err := tx.Model(&Event{}).Where(Event{
		ID: uuid.MustParse(t.EventID),
	}).Take(&event).Error; err != nil {
		return err
	}

	event.Registers += 1

	if err := tx.Model(&Event{}).Where(Event{
		ID: uuid.MustParse(t.EventID),
	}).Updates(event).Error; err != nil {
		return err
	}

	return nil
}
