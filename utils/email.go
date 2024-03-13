package utils

import (
	"github.com/TEDxITS/website-backend-2024/config"

	"gopkg.in/gomail.v2"
)

type Email struct {
	Email   string
	Subject string
	Body    string
}

func SendMail(mail Email) error {

	emailConfig := config.NewEmailConfig()
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", emailConfig.AuthEmail)
	mailer.SetHeader("To", mail.Email)
	mailer.SetHeader("Subject", mail.Subject)
	mailer.SetBody("text/html", mail.Body)

	dialer := gomail.NewDialer(
		emailConfig.Host,
		emailConfig.Port,
		emailConfig.AuthEmail,
		emailConfig.AuthPassword,
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		return err
	}

	return nil
}
