package entity

import (
	"github.com/TEDxITS/website-backend-2024/constants"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	PE2RSVP struct {
		ID         uuid.UUID `json:"id" form:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
		Name       string    `json:"name" form:"name"`
		Email      string    `json:"email" form:"email"`
		Institute  string    `json:"institute" form:"institute"`
		Department string    `json:"department" form:"department"`
		StudentID  string    `json:"student_id" form:"student_id"`
		Batch      string    `json:"batch" form:"batch"`

		WillingToCome        *bool  `json:"willing_to_come" form:"willing_to_come"`
		WillingToBeContacted *bool  `json:"willing_to_be_contacted" form:"willing_to_be_contacted"`
		Essay                string `json:"essay" form:"essay" gorm:"comment:How do you see Indonesia in the next 10 years due to the influence of its politics?"`
	}
)

func (e *PE2RSVP) AfterCreate(tx *gorm.DB) error {
	if !*e.WillingToCome {
		return nil
	}

	var event Event
	if err := tx.Model(&Event{}).Where(Event{
		Name: constants.PE2Name,
	}).Take(&event).Error; err != nil {
		return err
	}

	event.Registers += 1

	if err := tx.Model(&Event{}).Where(Event{
		Name: constants.PE2Name,
	}).Updates(event).Error; err != nil {
		return err
	}

	return nil
}
