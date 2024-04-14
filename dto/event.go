package dto

import (
	"errors"
	"time"
)

const (
	MESSAGE_FAILED_GET_EVENT = "failed get event"
	MESSAGE_EVENT_NOT_FOUND  = "event not found"

	MESSAGE_SUCCESS_GET_EVENT = "success get event"
)

var (
	ErrFailedToFetch = errors.New("failed to fetch event")
)

type (
	EventResponse struct {
		ID          string    `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Price       int       `json:"price"`
		Capacity    int       `json:"capacity,omitempty"`
		Registers   int       `json:"registers,omitempty"`
		StartDate   time.Time `json:"start_date"`
		EndDate     time.Time `json:"end_date"`
	}

	// EventsDetailResponse struct {
	// 	Events []EventDetailResponse `json:"events"`
	// }
)
