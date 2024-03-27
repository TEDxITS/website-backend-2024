package entity

import "github.com/google/uuid"

type Ticket struct {
	TicketID uuid.UUID `json:"ticket_id" gorm:"primary_key" form:"ticket_id"`
	UserID   string    `json:"user_id" form:"user_id"`
	EventID  string    `json:"event_id" form:"event_id"`

	User  User  `gorm:"foreignKey:UserID"`
	Event Event `gorm:"foreignKey:EventID"`

	Timestamp
}
