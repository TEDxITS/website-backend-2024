package entity

type Ticket struct {
	UserID  string `json:"user_id" form:"user_id"`
	EventID string `json:"event_id" form:"event_id"`

	User  User  `gorm:"foreignKey:UserID"`
	Event Event `gorm:"foreignKey:EventID"`

	Timestamp
}
