package repository

import (
	"github.com/TEDxITS/website-backend-2024/entity"

	"gorm.io/gorm"
)

type (
	TicketRepository interface {
		CreateTicket(ticket entity.Ticket) (entity.Ticket, error)
	}

	ticketRepository struct {
		db *gorm.DB
	}
)

func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepository{
		db: db,
	}
}

func (r *ticketRepository) CreateTicket(ticket entity.Ticket) (entity.Ticket, error) {
	err := r.db.Create(&ticket).Error
	if err != nil {
		return entity.Ticket{}, err
	}

	return ticket, nil
}
