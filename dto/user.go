package dto

import (
	"errors"
)

const (
	// Failed
	MESSAGE_FAILED_REGISTER_USER             = "failed create user"
	MESSAGE_FAILED_GET_USER                  = "failed get user"
	MESSAGE_FAILED_LOGIN                     = "failed login"
	MESSAGE_FAILED_UPDATE_USER               = "failed update user"
	MESSAGE_FAILED_VERIFY_USER               = "failed verify user"
	MESSAGE_FAILED_RESEND_VERIFY_EMAIL       = "failed resend verify email"
	MESSAGE_FAILED_RESET_PASSWORD            = "failed reset password"
	MESSAGE_FAILED_SEND_RESET_PASSWORD_EMAIL = "failed send reset password email"

	// Success
	MESSAGE_SUCCESS_REGISTER_USER             = "success create user. Please verify your email to activate your account"
	MESSAGE_SUCCESS_GET_USER                  = "success get user"
	MESSAGE_SUCCESS_LOGIN                     = "success login"
	MESSAGE_SUCCESS_UPDATE_USER               = "success update user"
	MESSAGE_SUCCESS_VERIFY_USER               = "success verify user"
	MESSAGE_SUCCESS_RESEND_VERIFY_EMAIL       = "success resend verify email"
	MESSAGE_SUCCESS_RESET_PASSWORD            = "success reset password"
	MESSAGE_SUCCESS_SEND_RESET_PASSWORD_EMAIL = "success send reset password email"
)

var (
	ErrRoleNotAllowed             = errors.New("denied access for \"%v\" role")
	ErrCreateUser                 = errors.New("failed to create user")
	ErrGetUserById                = errors.New("failed to get user by id")
	ErrEmailAlreadyExists         = errors.New("email already exist")
	ErrUpdateUser                 = errors.New("failed to update user")
	ErrUserNotFound               = errors.New("user not found")
	ErrCredentialsNotMatched      = errors.New("credentials not matched")
	ErrAccountNotVerified         = errors.New("account not verified")
	ErrEmailFormatInvalid         = errors.New("email format invalid")
	ErrAccountAlreadyVerified     = errors.New("account already verified")
	ErrGenerateVerificationEmail  = errors.New("failed to generate verification email")
	ErrSendEmail                  = errors.New("failed to send email")
	ErrDecryptToken               = errors.New("failed to decrypt token")
	ErrVerifyEmail                = errors.New("failed to verify email")
	ErrInvalidToken               = errors.New("invalid token")
	ErrHashPassword               = errors.New("failed to hash password")
	ErrGenerateResetPasswordEmail = errors.New("failed to generate reset password email")
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

	UserResendVerifyEmailRequest struct {
		Email string `json:"email" form:"email" binding:"required"`
	}

	UserSendResetPassworRequest struct {
		Email string `json:"email" form:"email" binding:"required"`
	}

	UserResendVerifyEmailResponse struct {
		Email string `json:"email"`
	}

	UserResetPasswordRequest struct {
		Password string `json:"password" form:"password" binding:"required"`
	}
)
