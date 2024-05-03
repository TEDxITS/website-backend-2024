package dto

import (
	"errors"
	"mime/multipart"
	"time"
)

const (
	// failed
	MESSAGE_FAILED_CREATE_TICKET   = "failed create ticket"
	MESSAGE_FAILED_GET_TICKET      = "failed get ticket"
	MESSAGE_FAILED_CONFIRM_PAYMENT = "failed confirm payment"
	MESSAGE_FAILED_CHECK_IN        = "failed check in"

	// success
	MESSAGE_SUCCESS_CREATE_TICKET   = "success create ticket"
	MESSAGE_SUCCESS_GET_TICKET      = "success get ticket"
	MESSAGE_SUCCESS_CONFIRM_PAYMENT = "success confirm payment"
	MESSAGE_SUCCESS_CHECK_IN        = "success check in"

	MAIN_EVENT_CLOSED = "closed"
	MAIN_EVENT_OPEN   = "open"
	MAIN_EVENT_FULL   = "full"
)

var (
	ErrUserNotInTransaction     = errors.New("user not in transaction")
	ErrFailedToStorePaymentFile = errors.New("failed to store payment file")
	ErrMainEventFull            = errors.New("main event full")
	ErrMainEventNotYetOpen      = errors.New("main event not yet open")
	ErrMainEventClosed          = errors.New("main event closed")
	ErrMismatchData             = errors.New("mismatch data")
	ErrOpeningPaymentFile       = errors.New("failed to open payment file")
	ErrFailedToDownloadFile     = errors.New("failed to download file")
	ErrMaxFileSize5MB           = errors.New("max file size is 5MB")
	ErrFileMustBeImage          = errors.New("file must be an image (jpg/jpeg/png)")
	ErrFileNotFound             = errors.New("file not found")
)

type (
	MainEventConfirmPaymentRequest struct {
		Code string `json:"code" form:"code" binding:"required"`
	}

	MainEventCheckInRequest struct {
		Code string `json:"code" form:"code" binding:"required"`
	}

	MainEventStatusResponse struct {
		EarlyBird MainEventStatusDetail `json:"early_bird"`
		PreSale   MainEventStatusDetail `json:"pre_sale"`
		Normal    MainEventStatusDetail `json:"normal"`
	}

	MainEventStatusDetail struct {
		Status      string        `json:"status"`
		UntilOpen   RemainingTime `json:"until_open"`
		UntilClosed RemainingTime `json:"until_closed"`
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

	MainEventRegister struct {
		EventID     string                `json:"event_id" form:"event_id"`
		Handphone   string                `json:"handphone" form:"handphone"`
		Birthdate   time.Time             `json:"birthdate" form:"birthdate"`
		Seat        string                `json:"seat" form:"seat"`
		PaymentFile *multipart.FileHeader `json:"payment_file" form:"payment_file"`
	}
)
