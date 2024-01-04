package utils

import (
	"net/mail"
	"strings"
)

func ValidateLink(link string) bool {
	return strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://")
}

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
