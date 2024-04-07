package dto

import "errors"

const (
	MESSAGE_FAILED_GET_DATA_FROM_BODY = "failed get data from body"
	MESSAGE_FAILED_VERIFY_TOKEN       = "failed to verify JWT token"
)

var (
	ErrTokenInvalid  = errors.New("token invalid")
	ErrTokenExpired  = errors.New("token expired")
	ErrTokenNotFound = errors.New("token not found")
)
