package repository

import (
	"math"

	"github.com/TEDxITS/website-backend-2024/constants"
	"github.com/TEDxITS/website-backend-2024/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type (
	TicketRepository interface {
		CreateTicket(ticket entity.Ticket) (entity.Ticket, error)
		JoinGetAllPaginationME(search string, limit, page int) ([]entity.Ticket, int64, int64, error)
		JoinGetAllPaginationPE3(search string, limit, page int) ([]entity.Ticket, int64, int64, error)
		FindByUserID(userID string) (entity.Ticket, error)
		UpdateTicket(ticket entity.Ticket) (entity.Ticket, error)
		GetTicketByUserId(userId string) (entity.Ticket, error)
		FindByTicketID(ticketID string) (entity.Ticket, error)
		GetTicketById(id string) (entity.Ticket, error)
		CountME() (int64, int64, int64, error)
		CountPE3() (int64, int64, int64, error)
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

func (r *ticketRepository) JoinGetAllPaginationME(search string, limit, page int) ([]entity.Ticket, int64, int64, error) {
	var tickets []entity.Ticket
	var count int64

	if search != "" {
		err := r.db.
			Model(&entity.Ticket{}).
			Joins("JOIN users ON tickets.user_id = users.id").
			Joins("JOIN events ON tickets.event_id = events.id").
			Where("users.name LIKE ?", "%"+search+"%").
			Where("event_id <> ?", constants.PreEvent3ID).
			Count(&count).Error
		if err != nil {
			return nil, 0, 0, err
		}
	} else {
		err := r.db.Model(&entity.Ticket{}).Where("event_id <> ?", constants.PreEvent3ID).Count(&count).Error
		if err != nil {
			return nil, 0, 0, err
		}
	}

	maxPage := int64(math.Ceil(float64(count) / float64(limit)))
	offset := (page - 1) * limit

	err := r.db.
		Model(&entity.Ticket{}).
		Joins("JOIN users ON tickets.user_id = users.id").
		Joins("JOIN events ON tickets.event_id = events.id").
		Preload(clause.Associations).
		Where("users.name LIKE ?", "%"+search+"%").
		Where("event_id <> ?", constants.PreEvent3ID).
		Offset(offset).
		Limit(limit).
		Find(&tickets).Error
	if err != nil {
		return nil, 0, 0, err
	}

	return tickets, maxPage, count, nil
}

func (r *ticketRepository) JoinGetAllPaginationPE3(search string, limit, page int) ([]entity.Ticket, int64, int64, error) {
	var tickets []entity.Ticket
	var count int64

	if search != "" {
		err := r.db.
			Model(&entity.Ticket{}).
			Joins("JOIN users ON tickets.user_id = users.id").
			Joins("JOIN events ON tickets.event_id = events.id").
			Where("users.name LIKE ?", "%"+search+"%").
			Where("event_id = ?", constants.PreEvent3ID).
			Count(&count).Error
		if err != nil {
			return nil, 0, 0, err
		}
	} else {
		err := r.db.Model(&entity.Ticket{}).Where("event_id = ?", constants.PreEvent3ID).Count(&count).Error
		if err != nil {
			return nil, 0, 0, err
		}
	}

	maxPage := int64(math.Ceil(float64(count) / float64(limit)))
	offset := (page - 1) * limit

	err := r.db.
		Model(&entity.Ticket{}).
		Joins("JOIN users ON tickets.user_id = users.id").
		Joins("JOIN events ON tickets.event_id = events.id").
		Preload(clause.Associations).
		Where("users.name LIKE ?", "%"+search+"%").
		Where("event_id = ?", constants.PreEvent3ID).
		Offset(offset).
		Limit(limit).
		Find(&tickets).Error
	if err != nil {
		return nil, 0, 0, err
	}

	return tickets, maxPage, count, nil
}

func (r *ticketRepository) GetTicketByUserId(userId string) (entity.Ticket, error) {
	var ticket entity.Ticket
	if err := r.db.Where("user_id = ?", userId).Take(&ticket).Error; err != nil {
		return entity.Ticket{}, err
	}
	return ticket, nil
}

func (r *ticketRepository) GetTicketById(id string) (entity.Ticket, error) {
	var ticket entity.Ticket
	if err := r.db.Where("ticket_id = ?", id).Take(&ticket).Error; err != nil {
		return entity.Ticket{}, err
	}
	return ticket, nil
}

func (r *ticketRepository) CountME() (int64, int64, int64, error) {
	var total int64
	if err := r.db.Model(&entity.Ticket{}).Where("event_id <> ?", constants.PreEvent3ID).Count(&total).Error; err != nil {
		return 0, 0, 0, err
	}

	var confirmed int64
	if err := r.db.Model(&entity.Ticket{}).Where("payment_confirmed = ? AND event_id <> ?", true, constants.PreEvent3ID).Count(&confirmed).Error; err != nil {
		return 0, 0, 0, err
	}

	var checked int64
	if err := r.db.Model(&entity.Ticket{}).Where("checked_in = ? AND event_id <> ?", true, constants.PreEvent3ID).Count(&checked).Error; err != nil {
		return 0, 0, 0, err
	}

	return total, confirmed, checked, nil
}

func (r *ticketRepository) CountPE3() (int64, int64, int64, error) {
	var total int64
	if err := r.db.Model(&entity.Ticket{}).Where("event_id = ?", constants.PreEvent3ID).Count(&total).Error; err != nil {
		return 0, 0, 0, err
	}

	var confirmed int64
	if err := r.db.Model(&entity.Ticket{}).Where("payment_confirmed = ? AND event_id = ?", true, constants.PreEvent3ID).Count(&confirmed).Error; err != nil {
		return 0, 0, 0, err
	}

	var checked int64
	if err := r.db.Model(&entity.Ticket{}).Where("checked_in = ? AND event_id = ?", true, constants.PreEvent3ID).Count(&checked).Error; err != nil {
		return 0, 0, 0, err
	}

	return total, confirmed, checked, nil
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
