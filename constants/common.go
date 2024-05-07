package constants

import "time"

var BASE_URL string

const (
	ENUM_ROLE_ADMIN = "admin"
	ENUM_ROLE_USER  = "user"

	ENUM_RUN_PRODUCTION  = "production"
	ENUM_RUN_DEVELOPMENT = "development"

	ENUM_PAGINATION_LIMIT = 10
	ENUM_PAGINATION_PAGE  = 1

	CTX_KEY_TOKEN     = "TOKEN"
	CTX_KEY_USER_ID   = "user_id"
	CTX_KEY_ROLE_NAME = "role"

	JWT_EXPIRE_TIME_IN_MINUTES = 120

	WSOCKET_AUTH_TIME_LIMIT        = time.Second * time.Duration(10)
	WSOCKET_TRANSACTION_TIME_LIMIT = (time.Minute * time.Duration(3)) + (time.Second * time.Duration(10))
)

var (
	CORS_ALLOWED_ORIGIN = []string{"http://localhost:3000", "https://tedxits.com", "http://tedxits.com", "https://www.tedxits.com", "http://www.tedxits.com"}
)
