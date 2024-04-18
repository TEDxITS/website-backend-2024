package dto

import "time"

type (
	MainEventConfirmPaymentRequest struct {
		Email string `json:"email" form:"email" binding:"required"`
	}

	MainEventCheckInRequest struct {
		Code string `json:"code" form:"code" binding:"required"`
	}

	MainEventDetailResponse struct {
		EventResponse
		Status        bool          `json:"status"`
		RemainingTime RemainingTime `json:"remaining_time"`
	}

	RemainingTime struct {
		Days    int `json:"days"`
		Hours   int `json:"hours"`
		Minutes int `json:"minutes"`
		Seconds int `json:"seconds"`
	}

	MainEventPaginationData struct {
		ID        string `json:"id" form:"id"`
		Name      string `json:"name" form:"name"`
		Email     string `json:"email" form:"email"`
		Confirmed bool   `json:"confirmed" form:"confirmed"`
		CheckedIn bool   `json:"checked_in" form:"checked_in"`
		EventName string `json:"event_name" form:"event_name"`
		Price     int    `json:"price" form:"price"`
	}

	MainEventPaginationResponse struct {
		Data []MainEventPaginationData `json:"data"`
		PaginationMetadata
	}

	MainEventResponse struct {
		ID        string `json:"id" form:"id"`
		Name      string `json:"name" form:"name"`
		Email     string `json:"email" form:"email"`
		Confirmed bool   `json:"confirmed" form:"confirmed"`
		CheckedIn bool   `json:"checked_in" form:"checked_in"`
		EventName string `json:"event_name" form:"event_name"`
		Price     int    `json:"price" form:"price"`

		Handphone    string    `json:"handphone" form:"handphone"`
		Birthdate    time.Time `json:"birthdate" form:"birthdate"`
		Seat         string    `json:"seat" form:"seat"`
		Payment      string    `json:"payment" form:"payment"`
		WithKit      bool      `json:"with_kit" form:"with_kit"`
		RegisterDate time.Time `json:"registerdate" form:"registerdate"`
	}

	MainEventCounter struct {
		Total             int64 `json:"total" form:"total"`
		ConfirmedPayments int64 `json:"confirmed_payments" form:"confirmed_payments"`
		CheckedIns        int64 `json:"checked_ins" form:"checked_ins"`
	}
)
