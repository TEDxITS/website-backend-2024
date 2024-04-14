package dto

import (
	"errors"

	"github.com/google/uuid"
)

const (
	// Failed
	MESSAGE_FAILED_CREATE_TICKET = "failed create ticket"
	MESSAGE_FAILED_GET_TICKET    = "failed get ticket"

	// Success
	MESSAGE_SUCCESS_CREATE_TICKET = "success create ticket"
	MESSAGE_SUCCESS_GET_TICKET    = "success get ticket"
)

var (
	ErrCreateTicket           = errors.New("failed to create ticket")
	ErrPE2RSVPNotOpen         = errors.New("pre event 2 RSVP is not yet open")
	ErrPE2RSVPClosed          = errors.New("pre event 2 RSVP is closed")
	ErrPE2RSVPFull            = errors.New("pre event 2 RSVP is full")
	ErrPE2RSVPEmailRegistered = errors.New("email already registered")
)

type (
	TicketPE2RSVPRequest struct {
		Name       string `json:"name" form:"name" binding:"required"`
		Email      string `json:"email" form:"email" binding:"required"`
		Institute  string `json:"institute" form:"institute" binding:"required"`
		Department string `json:"department" form:"department"`
		StudentID  string `json:"student_id" form:"student_id"`
		Batch      string `json:"batch" form:"batch"`

		WillingToCome        bool   `json:"willing_to_come" form:"willing_to_come" binding:"boolean"`
		WillingToBeContacted bool   `json:"willing_to_be_contacted" form:"willing_to_be_contacted" binding:"boolean"`
		Essay                string `json:"essay" form:"essay" binding:"required"`
	}

	TicketPE2RSVPResponse struct {
		ID         uuid.UUID `json:"id" form:"id"`
		Name       string    `json:"name" form:"name"`
		Email      string    `json:"email" form:"email"`
		Institute  string    `json:"institute" form:"institute"`
		Department string    `json:"department" form:"department"`
		StudentID  string    `json:"student_id" form:"student_id"`
		Batch      string    `json:"batch" form:"batch"`

		WillingToCome        bool   `json:"willing_to_come" form:"willing_to_come"`
		WillingToBeContacted bool   `json:"willing_to_be_contacted" form:"willing_to_be_contacted"`
		Essay                string `json:"essay" form:"essay"`
	}

	TicketPE2RSVPCounter struct {
		Total   int64 `json:"total" form:"total"`
		Attends int64 `json:"attends" form:"attends"`
	}

	TicketPE2RSVPPaginationData struct {
		ID                   uuid.UUID `json:"id" form:"id"`
		Name                 string    `json:"name" form:"name"`
		Institute            string    `json:"institute" form:"institute"`
		Batch                string    `json:"batch" form:"batch"`
		WillingToCome        bool      `json:"willing_to_come" form:"willing_to_come"`
		WillingToBeContacted bool      `json:"willing_to_be_contacted" form:"willing_to_be_contacted"`
	}

	TicketPE2RSVPPaginationResponse struct {
		Data []TicketPE2RSVPPaginationData `json:"data"`
		PaginationMetadata
	}
)
