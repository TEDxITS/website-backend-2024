package entity

import "github.com/google/uuid"

type (
	PreEventTwoRSVP struct {
		ID         uuid.UUID `json:"id" form:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
		Name       string    `json:"name" form:"name"`
		Email      string    `json:"email" form:"email"`
		Institute  string    `json:"institute" form:"institute"`
		Department string    `json:"department" form:"department"`
		StudentID  string    `json:"student_id" form:"student_id"`
		Batch      string    `json:"batch" form:"batch"`

		WillingToCome        bool   `json:"willing_to_come" form:"willing_to_come"`
		WillingToBeContacted bool   `json:"willing_to_be_contacted" form:"willing_to_be_contacted"`
		Essay                string `json:"essay" form:"essay" gorm:"comment:How do you see Indonesia in the next 10 years due to the influence of its politics?"`
	}
)
