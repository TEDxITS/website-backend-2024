package repository

import (
	"math"

	"github.com/TEDxITS/website-backend-2024/entity"

	"gorm.io/gorm"
)

type (
	TicketRepository interface {
		CreateTicket(ticket entity.Ticket) (entity.Ticket, error)
		GetAllPagination(search string, limit, page int) ([]entity.User, int64, int64, error)
		GetTicketByUserId(userId string) (entity.Ticket, error)
		GetEventByID(id string) (entity.Event, error)
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

func (r *ticketRepository) GetAllPagination(search string, limit, page int) ([]entity.User, int64, int64, error) {
	var rsvps []entity.User
	var count int64

	if search != "" {
		err := r.db.Model(&entity.User{}).
			Where("name LIKE ?", "%"+search+"%").
			Count(&count).Error
		if err != nil {
			return nil, 0, 0, err
		}
	} else {
		err := r.db.Model(&entity.User{}).Count(&count).Error
		if err != nil {
			return nil, 0, 0, err
		}
	}

	maxPage := int64(math.Ceil(float64(count) / float64(limit)))
	offset := (page - 1) * limit

	err := r.db.Model(&entity.User{}).
		Where("name LIKE ?", "%"+search+"%").
		Offset(offset).Limit(limit).Find(&rsvps).Error
	if err != nil {
		return nil, 0, 0, err
	}

	return rsvps, maxPage, count, nil
}

func (r *ticketRepository) GetTicketByUserId(userId string) (entity.Ticket, error) {
	var ticket entity.Ticket
	if err := r.db.Where("user_id = ?", userId).Take(&ticket).Error; err != nil {
		return entity.Ticket{}, err
	}
	return ticket, nil
}

func (r *ticketRepository) GetEventByID(id string) (entity.Event, error) {
	var event entity.Event
	if err := r.db.Where("id = ?", id).Take(&event).Error; err != nil {
		return entity.Event{}, err
	}
	return event, nil
}
