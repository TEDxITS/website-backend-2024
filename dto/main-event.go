package dto

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
)
