package dto

import "errors"

const (
	MESSAGE_FAILED_CREATE_LINK_SHORTENER = "failed create link shortener"
	MESSAGE_FAILED_GET_LINK_SHORTENER    = "failed get link shortener"

	MESSAGE_SUCCESS_CREATE_LINK_SHORTENER = "success create link shortener"
	MESSAGE_SUCCESS_GET_LINK_SHORTENER    = "success get link shortener"
)

var (
	ErrLinkFormatInvalid   = errors.New("link format should start with http:// or https://")
	ErrCreateLinkShortener = errors.New("failed to create link shortener")
	ErrAliasHasBeenTaken   = errors.New("alias has been taken")
	ErrAliasNotFound       = errors.New("alias not found")
)

type (
	LinkShortenerRequest struct {
		Alias string `json:"alias" form:"alias" binding:"required"`
		Link  string `json:"link" form:"link" binding:"required"`
	}

	LinkShortenerResponse struct {
		ID    string `json:"id,omitempty"`
		Alias string `json:"alias"`
		Link  string `json:"link"`
	}

	LinkShortenerPaginationResponse struct {
		Data []LinkShortenerResponse `json:"data"`
		PaginationMetadata
	}
)
