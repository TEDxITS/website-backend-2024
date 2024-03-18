package dto

import (
	"errors"
)

const (
	// Failed
	MESSAGE_FAILED_REGISTER_USER = "failed create user"
	MESSAGE_FAILED_GET_USER      = "failed get user"
	MESSAGE_FAILED_LOGIN         = "failed login"
	MESSAGE_FAILED_UPDATE_USER   = "failed update user"
	MESSAGE_FAILED_VERIFY_USER   = "failed verify user"

	// Success
	MESSAGE_SUCCESS_REGISTER_USER = "success create user"
	MESSAGE_SUCCESS_GET_USER      = "success get user"
	MESSAGE_SUCCESS_LOGIN         = "success login"
	MESSAGE_SUCCESS_UPDATE_USER   = "success update user"
	MESSAGE_SUCCESS_VERIFY_USER   = "success verify user"
)

var (
	ErrRoleNotAllowed        = errors.New("denied access for \"%v\" role")
	ErrCreateUser            = errors.New("failed to create user")
	ErrGetUserById           = errors.New("failed to get user by id")
	ErrEmailAlreadyExists    = errors.New("email already exist")
	ErrUpdateUser            = errors.New("failed to update user")
	ErrUserNotFound          = errors.New("user not found")
	ErrCredentialsNotMatched = errors.New("credentials not matched")
	ErrAccountNotVerified    = errors.New("account not verified")
	ErrEmailFormatInvalid    = errors.New("email format invalid")
)

type (
	UserRequest struct {
		Name     string `json:"name" form:"name" binding:"required"`
		Email    string `json:"email" form:"email" binding:"required"`
		Password string `json:"password" form:"password" binding:"required"`
	}

	UserLoginRequest struct {
		Email    string `json:"email" form:"email" binding:"required"`
		Password string `json:"password" form:"password" binding:"required"`
	}

	UserResponse struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		Email      string `json:"email"`
		RoleID     string `json:"role_id,omitempty"`
		Role       string `json:"role,omitempty"`
		IsVerified bool   `json:"is_verified"`
	}

	UserPaginationResponse struct {
		Data []UserResponse `json:"data"`
		PaginationMetadata
	}
)
