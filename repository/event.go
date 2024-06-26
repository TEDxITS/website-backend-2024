package repository

import (
	"github.com/TEDxITS/website-backend-2024/constants"
	"github.com/TEDxITS/website-backend-2024/entity"
	"gorm.io/gorm"
)

type (
	EventRepository interface {
		GetAll() ([]entity.Event, error)
		GetByID(string) (entity.Event, error)
		GetPE2Detail() (entity.Event, error)
		GetAllExcept(eventID string) ([]entity.Event, error)
	}

	eventRepository struct {
		db *gorm.DB
	}
)

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{
		db: db,
	}
}

func (r *eventRepository) GetAll() ([]entity.Event, error) {
	var events []entity.Event
	if err := r.db.Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

func (r *eventRepository) GetByID(id string) (entity.Event, error) {
	var event entity.Event
	if err := r.db.Where("id = ?", id).Take(&event).Error; err != nil {
		return entity.Event{}, err
	}
	return event, nil
}

func (r *eventRepository) GetPE2Detail() (entity.Event, error) {
	var event entity.Event
	if err := r.db.Where("name = ?", constants.PE2Name).Take(&event).Error; err != nil {
		return entity.Event{}, err
	}
	return event, nil
}

func (r *eventRepository) GetAllExcept(eventID string) ([]entity.Event, error) {
	var events []entity.Event
	if err := r.db.Where("id != ?", eventID).Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}
