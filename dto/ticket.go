package dto

import (
	"errors"
)

const (
	// Failed
	MESSAGE_FAILED_CREATE_TICKET = "failed create user"

	// Success
	MESSAGE_SUCCESS_CREATE_TICKET = "success create ticket"
)

var (
	ErrCreateTicket = errors.New("failed to create ticket")
)

type (
	TicketRequest struct {
		UserID  string `json:"user_id" form:"user_id" binding:"required"`
		EventID string `json:"event_id" form:"event_id" binding:"required"`
	}

	TicketResponse struct {
		UserID  string `json:"user_id" form:"user_id" binding:"required"`
		EventID string `json:"event_id" form:"event_id" binding:"required"`
	}
)
