package entity

import "time"

type Ticket struct {
	TicketID string `json:"ticket_id" gorm:"primary_key" form:"ticket_id"`
	UserID   string `json:"user_id" form:"user_id"`
	EventID  string `json:"event_id" form:"event_id"`

	Handphone string    `json:"handphone" form:"handphone"`
	Birthdate time.Time `json:"birthdate" form:"birthdate"`
	Seat      string    `json:"seat" form:"seat"`
	Payment   string    `json:"payment" form:"payment"`

	User  User  `gorm:"foreignKey:UserID"`
	Event Event `gorm:"foreignKey:EventID"`

	Timestamp
}
