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
		FindByUserID(userID string) (entity.Ticket, error)
		UpdateTicket(ticket entity.Ticket) (entity.Ticket, error)
		GetTicketByUserId(userId string) (entity.Ticket, error)
		FindByTicketID(ticketID string) (entity.Ticket, error)
		GetEventById(id string) (entity.Event, error)
		GetTicketById(id string) (entity.Ticket, error)
		GetUserById(id string) (entity.User, error)
		CountTotal() (int64, error)
		CountConfirmedPayments() (int64, error)
		CountCheckedIns() (int64, error)
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

func (r *ticketRepository) GetEventById(id string) (entity.Event, error) {
	var event entity.Event
	if err := r.db.Where("id = ?", id).Take(&event).Error; err != nil {
		return entity.Event{}, err
	}
	return event, nil
}

func (r *ticketRepository) GetTicketById(id string) (entity.Ticket, error) {
	var ticket entity.Ticket
	if err := r.db.Where("ticket_id = ?", id).Take(&ticket).Error; err != nil {
		return entity.Ticket{}, err
	}
	return ticket, nil
}

func (r *ticketRepository) GetUserById(id string) (entity.User, error) {
	var user entity.User
	if err := r.db.Where("id = ?", id).Take(&user).Error; err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (r *ticketRepository) CountTotal() (int64, error) {
	var count int64
	if err := r.db.Model(&entity.Ticket{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ticketRepository) CountConfirmedPayments() (int64, error) {
	var count int64
	if err := r.db.Model(&entity.Ticket{}).Where("payment_confirmed = ?", true).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ticketRepository) CountCheckedIns() (int64, error) {
	var count int64
	if err := r.db.Model(&entity.Ticket{}).Where("checked_in = ?", true).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
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
