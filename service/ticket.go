package service

import (
	"context"

	"github.com/TEDxITS/website-backend-2024/dto"
	"github.com/TEDxITS/website-backend-2024/entity"
	"github.com/TEDxITS/website-backend-2024/repository"
)

type (
	TicketService interface {
		CreateTicket(ctx context.Context, req dto.TicketRequest) (dto.TicketResponse, error)
	}

	ticketService struct {
		ticketRepo repository.TicketRepository
	}
)

func NewTicketService(tr repository.TicketRepository) TicketService {
	return &ticketService{
		ticketRepo: tr,
	}
}

func (t *ticketService) CreateTicket(ctx context.Context, req dto.TicketRequest) (dto.TicketResponse, error) {
	ticket := entity.Ticket{
		UserID:  req.UserID,
		EventID: req.EventID,
	}

	res, err := t.ticketRepo.CreateTicket(ticket)
	if err != nil {
		return dto.TicketResponse{}, dto.ErrCreateTicket
	}

	return dto.TicketResponse{
		UserID:  res.UserID,
		EventID: res.EventID,
	}, nil
}
