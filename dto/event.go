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
	AdminEventResponse struct {
		ID          string    `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Price       int       `json:"price"`
		Capacity    int       `json:"capacity"`
		Registers   int       `json:"registers"`
		StartDate   time.Time `json:"start_date"`
		EndDate     time.Time `json:"end_date"`
	}

	UserEventResponse struct {
		ID          string    `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Price       int       `json:"price"`
		StartDate   time.Time `json:"start_date"`
		EndDate     time.Time `json:"end_date"`
	}

	EventResponse struct {
		Data interface{} `json:"data"`
	}
)
