package service

import (
	"bytes"
	"context"
	"os"
	"text/template"
	"time"

	"github.com/TEDxITS/website-backend-2024/dto"
	"github.com/TEDxITS/website-backend-2024/repository"
	"github.com/TEDxITS/website-backend-2024/utils"
)

type (
	MainEventService interface {
		ConfirmPayment(context.Context, dto.MainEventConfirmPaymentRequest) error
		CheckIn(context.Context, dto.MainEventCheckInRequest) error
		GetStatus(context.Context) ([]dto.MainEventDetailResponse, error)
	}

	mainEventService struct {
		eventRepo  repository.EventRepository
		userRepo   repository.UserRepository
		ticketRepo repository.TicketRepository
	}
)

func NewMainEventService(uRepo repository.UserRepository, tRepo repository.TicketRepository, eRepo repository.EventRepository) MainEventService {
	return &mainEventService{
		eventRepo:  eRepo,
		userRepo:   uRepo,
		ticketRepo: tRepo,
	}
}

func (s *mainEventService) ConfirmPayment(ctx context.Context, req dto.MainEventConfirmPaymentRequest) error {
	email, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return dto.ErrUserNotFound
	}

	ticket, err := s.ticketRepo.FindByUserID(email.ID.String())
	if err != nil {
		return dto.ErrTicketNotFound
	}

	confirmed := true
	ticket.PaymentConfirmed = &confirmed
	_, err = s.ticketRepo.UpdateTicket(ticket)
	if err != nil {
		return err
	}

	readHtml, err := os.ReadFile("./utils/template/confirmation_payment.html")

	if err != nil {
		return err
	}

	data := struct {
		Email    string
		TicketID string
	}{
		Email:    req.Email,
		TicketID: ticket.TicketID,
	}

	tmpl, err := template.New("custom").Parse(string(readHtml))
	if err != nil {
		return err
	}

	var strMail bytes.Buffer
	if err := tmpl.Execute(&strMail, data); err != nil {
		return err
	}

	emailData := utils.Email{
		Email:   req.Email,
		Subject: "Confirmation Payment",
		Body:    strMail.String(),
	}

	err = utils.SendMail(emailData)
	if err != nil {
		return dto.ErrSendEmail
	}

	return nil
}

func (s *mainEventService) CheckIn(ctx context.Context, req dto.MainEventCheckInRequest) error {
	ticket, err := s.ticketRepo.FindByTicketID(req.Code)
	if err != nil {
		return dto.ErrTicketNotFound
	}

	checked := true
	ticket.CheckedIn = &checked
	_, err = s.ticketRepo.UpdateTicket(ticket)
	if err != nil {
		return err
	}

	return nil
}

func (s *mainEventService) GetStatus(ctx context.Context) ([]dto.MainEventDetailResponse, error) {
	events, err := s.eventRepo.GetAllExcept("7de24efe-0aec-469a-bf0c-8fa8cae3ff3f")
	if err != nil {
		return []dto.MainEventDetailResponse{}, err
	}

	var result []dto.MainEventDetailResponse
	for _, e := range events {
		eventResponse := dto.MainEventDetailResponse{
			EventResponse: dto.EventResponse{
				ID:        e.ID.String(),
				Name:      e.Name,
				Price:     e.Price,
				StartDate: e.StartDate,
				EndDate:   e.EndDate,
			},
		}

		eventResponse.Status = true

		if e.Registers >= e.Capacity {
			eventResponse.Status = false
		}

		if time.Now().Before(e.StartDate.Add(-7 * time.Hour)) {
			eventResponse.Status = false
		}

		if time.Now().After(e.EndDate.Add(-7 * time.Hour)) {
			eventResponse.Status = false
		}

		difference := e.EndDate.Add(-7 * time.Hour).Sub(time.Now())

		total := int(difference.Seconds())
		days := int(total / (60 * 60 * 24))
		hours := int(total / (60 * 60) % 24)
		minutes := int(total/60) % 60
		seconds := int(total % 60)

		eventResponse.RemainingTime = dto.RemainingTime{
			Days:    days,
			Hours:   hours,
			Minutes: minutes,
			Seconds: seconds,
		}

		result = append(result, eventResponse)
	}

	return result, nil
}
