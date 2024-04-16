package service

import (
	"bytes"
	"context"
	"os"
	"text/template"
	"time"

	"github.com/TEDxITS/website-backend-2024/constants"
	"github.com/TEDxITS/website-backend-2024/dto"
	"github.com/TEDxITS/website-backend-2024/repository"
	"github.com/TEDxITS/website-backend-2024/utils"
)

type (
	MainEventService interface {
		ConfirmPayment(context.Context, dto.MainEventConfirmPaymentRequest) error
		CheckIn(context.Context, dto.MainEventCheckInRequest) error
		GetStatus(context.Context) ([]dto.MainEventDetailResponse, error)
		GetMainEventPaginated(context.Context, dto.PaginationQuery) (dto.TicketMainEventPaginationResponse, error)
		GetMainEventDetail(context.Context, string) (dto.TicketMainEventResponse, error)
		GetMainEventCounter(context.Context) (dto.TicketMainEventCounter, error)
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

func (s *mainEventService) GetMainEventPaginated(ctx context.Context, req dto.PaginationQuery) (dto.TicketMainEventPaginationResponse, error) {
	var limit int
	var page int

	limit = req.PerPage
	if limit <= 0 {
		limit = constants.ENUM_PAGINATION_LIMIT
	}

	page = req.Page
	if page <= 0 {
		page = constants.ENUM_PAGINATION_PAGE
	}

	rsvps, maxPage, count, err := s.ticketRepo.GetAllPagination(req.Search, limit, page)
	if err != nil {
		return dto.TicketMainEventPaginationResponse{}, err
	}

	var result []dto.TicketMainEventPaginationData
	for _, rsvp := range rsvps {
		ticket, _ := s.ticketRepo.GetTicketByUserId(rsvp.ID.String())
		event, _ := s.ticketRepo.GetEventById(ticket.EventID)
		result = append(result, dto.TicketMainEventPaginationData{
			ID:        ticket.TicketID,
			Name:      rsvp.Name,
			Email:     rsvp.Email,
			Confirmed: *ticket.PaymentConfirmed,
			CheckedIn: *ticket.CheckedIn,
			EventName: event.Name,
			Price:     event.Price,
		})
	}

	return dto.TicketMainEventPaginationResponse{
		Data: result,
		PaginationMetadata: dto.PaginationMetadata{
			Page:    page,
			PerPage: limit,
			MaxPage: maxPage,
			Count:   count,
		},
	}, nil
}

func (s *mainEventService) GetMainEventDetail(ctx context.Context, id string) (dto.TicketMainEventResponse, error) {
	ticket, err := s.ticketRepo.GetTicketById(id)
	if err != nil {
		return dto.TicketMainEventResponse{}, err
	}
	event, err := s.ticketRepo.GetEventById(ticket.EventID)
	if err != nil {
		return dto.TicketMainEventResponse{}, err
	}
	user, err := s.ticketRepo.GetUserById(ticket.UserID)
	if err != nil {
		return dto.TicketMainEventResponse{}, err
	}
	return dto.TicketMainEventResponse{
		ID:        ticket.TicketID,
		Name:      user.Name,
		Email:     user.Email,
		Confirmed: *ticket.PaymentConfirmed,
		CheckedIn: *ticket.CheckedIn,
		EventName: event.Name,
		Price:     event.Price,

		Handphone: ticket.Handphone,
		Birthdate: ticket.Birthdate,
		Seat:      ticket.Seat,
		Payment:   ticket.Payment,
		WithKit:   *event.WithKit,
	}, nil
}

func (s *mainEventService) GetMainEventCounter(ctx context.Context) (dto.TicketMainEventCounter, error) {
	total, err := s.ticketRepo.CountTotal()
	if err != nil {
		return dto.TicketMainEventCounter{}, err
	}

	confirmed_payments, err := s.ticketRepo.CountConfirmedPayments()
	if err != nil {
		return dto.TicketMainEventCounter{}, err
	}

	checked_ins, err := s.ticketRepo.CountCheckedIns()
	if err != nil {
		return dto.TicketMainEventCounter{}, err
	}

	return dto.TicketMainEventCounter{
		Total:             total,
		ConfirmedPayments: confirmed_payments,
		CheckedIns:        checked_ins,
	}, nil
}
