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
		RemainingTime time.Duration `json:"remaining_time"`
	}
)
