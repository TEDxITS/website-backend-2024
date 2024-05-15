package dto

import (
	"errors"
	"mime/multipart"
	"time"
)

const (
	MESSAGE_FAILED_GET_PE3_STATUS  = "Failed to get pre-Event 3 status"
	MESSAGE_SUCCESS_GET_PE3_STATUS = "Successfully get pre-Event 3 status"
)

var (
	ErrPreEvent3NotYetOpen = errors.New("pre-Event 3 registration is not yet open")
	ErrPreEvent3Closed     = errors.New("pre-Event 3 registration is closed")
)

type (
	PE3RSVPRegister struct {
		Handphone   string                `json:"handphone" form:"handphone" binding:"required"`
		Birthdate   time.Time             `json:"birthdate" form:"birthdate" binding:"required"`
		PaymentFile *multipart.FileHeader `json:"payment_file" form:"payment_file" binding:"required"`
	}

	PE3RSVPStatus struct {
		Status *bool `json:"status"`
	}
)
