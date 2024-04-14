package repository

import (
	"github.com/TEDxITS/website-backend-2024/entity"

	"gorm.io/gorm"
)

type (
	TicketRepository interface {
		CreateTicket(ticket entity.Ticket) (entity.Ticket, error)
		FindByUserID(userID string) (entity.Ticket, error)
		UpdateTicket(ticket entity.Ticket) (entity.Ticket, error)
		FindByTicketID(ticketID string) (entity.Ticket, error)
		FindAll() ([]entity.Ticket, error)
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

func (r *ticketRepository) FindByUserID(userID string) (entity.Ticket, error) {
	var ticket entity.Ticket
	err := r.db.Where("user_id = ?", userID).First(&ticket).Error
	if err != nil {
		return entity.Ticket{}, err
	}

	return ticket, nil
}

func (r *ticketRepository) UpdateTicket(ticket entity.Ticket) (entity.Ticket, error) {
	err := r.db.Save(&ticket).Error
	if err != nil {
		return entity.Ticket{}, err
	}

	return ticket, nil
}

func (r *ticketRepository) FindByTicketID(ticketID string) (entity.Ticket, error) {
	var ticket entity.Ticket
	err := r.db.Where("ticket_id = ?", ticketID).First(&ticket).Error
	if err != nil {
		return entity.Ticket{}, err
	}

	return ticket, nil
}

func (r *ticketRepository) FindAll() ([]entity.Ticket, error) {
	var tickets []entity.Ticket
	err := r.db.Find(&tickets).Error
	if err != nil {
		return []entity.Ticket{}, err
	}

	return tickets, nil
}
