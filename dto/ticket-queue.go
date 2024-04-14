package dto

import "errors"

const (
	WSOCKET_AUTH_REQUEST            = "TOKEN %v"
	WSOCKET_ENUM_NO_MERCH_REQUEST   = "MERCH 0"
	WSOCKET_ENUM_WITH_MERCH_REQUEST = "MERCH 1"

	WSOCKET_PAYMENT_CODE        = "PAYMENT CODE %v"
	WSOCKET_QUEUE_NUMBER        = "WAITING QUEUE %v"
	WSOCKET_TICKET_AVAILABILITY = "M=%v | N=%v"

	WSOCKET_AUTH_SUCCESS        = "authentication successful"
	WSOCKET_TRANSACTION_START   = "proceed transaction"
	WSOCKET_TRANSACTION_SUCCESS = "transaction successful"
)

var (
	ErrWSAlreadyInQueue    = errors.New("already in queue")
	ErrWSBadRequest        = errors.New("bad request")
	ErrWSInvalidToken      = errors.New("invalid token")
	ErrWSInvalidCommand    = errors.New("invalid command")
	ErrWSMainEventFull     = errors.New("main event is full")
	ErrWSCommunicateWithDB = errors.New("error fetching data")
)
