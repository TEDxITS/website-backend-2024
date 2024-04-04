package service

import (
	"context"
	"time"

	"github.com/TEDxITS/website-backend-2024/constants"
	"github.com/TEDxITS/website-backend-2024/dto"
	"github.com/TEDxITS/website-backend-2024/entity"
	"github.com/TEDxITS/website-backend-2024/repository"
)

type (
	TicketService interface {
		CreatePE2RSVP(context.Context, dto.TicketPE2RSVPRequest) (dto.TicketPE2RSVPResponse, error)
		GetPE2RSVPPaginated(context.Context, dto.PaginationQuery) (dto.TicketPE2RSVPPaginationResponse, error)
		GetPE2RSVPDetail(context.Context, string) (dto.TicketPE2RSVPResponse, error)
		GetPE2RSVPCounter(context.Context) (dto.TicketPE2RSVPCounter, error)
	}

	ticketService struct {
		eventRepo   repository.EventRepository
		pe2RSVPRepo repository.PE2RSVPRepository
	}
)

func NewTicketService(eventRepo repository.EventRepository, pe2RSVPRepo repository.PE2RSVPRepository) TicketService {
	return &ticketService{
		eventRepo:   eventRepo,
		pe2RSVPRepo: pe2RSVPRepo,
	}
}

func (s *ticketService) CreatePE2RSVP(ctx context.Context, req dto.TicketPE2RSVPRequest) (dto.TicketPE2RSVPResponse, error) {
	event, err := s.eventRepo.GetPE2Detail()
	if err != nil {
		return dto.TicketPE2RSVPResponse{}, err
	}

	if event.Registers >= event.Capacity {
		return dto.TicketPE2RSVPResponse{}, dto.ErrPE2RSVPFull
	}

	if time.Now().Before(event.StartDate) {
		return dto.TicketPE2RSVPResponse{}, dto.ErrPE2RSVPNotOpen
	}

	if time.Now().After(event.EndDate) {
		return dto.TicketPE2RSVPResponse{}, dto.ErrPE2RSVPClosed
	}

	exist, err := s.pe2RSVPRepo.CheckEmailExist(req.Email)
	if err != nil {
		return dto.TicketPE2RSVPResponse{}, err
	}

	if exist {
		return dto.TicketPE2RSVPResponse{}, dto.ErrPE2RSVPEmailRegistered
	}

	rsvp := entity.PE2RSVP{
		Name:                 req.Name,
		Email:                req.Email,
		Institute:            req.Institute,
		Department:           req.Department,
		StudentID:            req.StudentID,
		Batch:                req.Batch,
		WillingToCome:        &req.WillingToCome,
		WillingToBeContacted: &req.WillingToBeContacted,
		Essay:                req.Essay,
	}

	res, err := s.pe2RSVPRepo.Create(rsvp)
	if err != nil {
		return dto.TicketPE2RSVPResponse{}, err
	}

	return dto.TicketPE2RSVPResponse{
		ID:                   res.ID,
		Name:                 res.Name,
		Email:                res.Email,
		Institute:            res.Institute,
		Department:           res.Department,
		StudentID:            res.StudentID,
		Batch:                res.Batch,
		WillingToCome:        *res.WillingToCome,
		WillingToBeContacted: *res.WillingToBeContacted,
		Essay:                res.Essay,
	}, nil
}

func (s *ticketService) GetPE2RSVPPaginated(ctx context.Context, req dto.PaginationQuery) (dto.TicketPE2RSVPPaginationResponse, error) {
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

	rsvps, maxPage, count, err := s.pe2RSVPRepo.GetAllPagination(req.Search, limit, page)
	if err != nil {
		return dto.TicketPE2RSVPPaginationResponse{}, err
	}

	var result []dto.TicketPE2RSVPPaginationData
	for _, rsvp := range rsvps {
		result = append(result, dto.TicketPE2RSVPPaginationData{
			ID:                   rsvp.ID,
			Name:                 rsvp.Name,
			Institute:            rsvp.Institute,
			Batch:                rsvp.Batch,
			WillingToCome:        *rsvp.WillingToCome,
			WillingToBeContacted: *rsvp.WillingToBeContacted,
		})
	}

	return dto.TicketPE2RSVPPaginationResponse{
		Data: result,
		PaginationMetadata: dto.PaginationMetadata{
			Page:    page,
			PerPage: limit,
			MaxPage: maxPage,
			Count:   count,
		},
	}, nil
}

func (s *ticketService) GetPE2RSVPDetail(ctx context.Context, id string) (dto.TicketPE2RSVPResponse, error) {
	attendee, err := s.pe2RSVPRepo.GetById(id)
	if err != nil {
		return dto.TicketPE2RSVPResponse{}, err
	}

	return dto.TicketPE2RSVPResponse{
		Name:                 attendee.Name,
		Email:                attendee.Email,
		Institute:            attendee.Institute,
		Department:           attendee.Department,
		StudentID:            attendee.StudentID,
		Batch:                attendee.Batch,
		WillingToCome:        *attendee.WillingToCome,
		WillingToBeContacted: *attendee.WillingToBeContacted,
		Essay:                attendee.Essay,
	}, nil
}

func (s *ticketService) GetPE2RSVPCounter(ctx context.Context) (dto.TicketPE2RSVPCounter, error) {
	total, err := s.pe2RSVPRepo.CountTotal()
	if err != nil {
		return dto.TicketPE2RSVPCounter{}, err
	}

	attends, err := s.pe2RSVPRepo.CountAttends()
	if err != nil {
		return dto.TicketPE2RSVPCounter{}, err
	}

	return dto.TicketPE2RSVPCounter{
		Total:   total,
		Attends: attends,
	}, nil
}
