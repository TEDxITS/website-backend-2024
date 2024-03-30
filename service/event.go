package service

import (
	"context"

	"github.com/TEDxITS/website-backend-2024/constants"
	"github.com/TEDxITS/website-backend-2024/dto"
	"github.com/TEDxITS/website-backend-2024/repository"
)

type (
	EventService interface {
		FindAll(ctx context.Context, userRole string) ([]dto.EventResponse, error)
		FindByID(ctx context.Context, id string, userRole string) (dto.EventResponse, error)
		GetPE2Detail(ctx context.Context, userRole string) (dto.EventResponse, error)
	}

	eventService struct {
		eventRepo repository.EventRepository
	}
)

func NewEventService(er repository.EventRepository) EventService {
	return &eventService{
		eventRepo: er,
	}
}

func (s *eventService) FindAll(ctx context.Context, userRole string) ([]dto.EventResponse, error) {
	events, err := s.eventRepo.GetAll()
	if err != nil {
		return []dto.EventResponse{}, err
	}

	var result []dto.EventResponse
	for _, event := range events {
		eventResponse := dto.EventResponse{
			ID:          event.ID.String(),
			Name:        event.Name,
			Description: event.Description,
			Price:       event.Price,
			StartDate:   event.StartDate,
			EndDate:     event.EndDate,
		}

		if userRole == constants.ENUM_ROLE_ADMIN {
			if event.Registers == 0 {
				event.Registers = 1
			}

			eventResponse.Capacity = event.Capacity
			eventResponse.Registers = event.Registers
		}

		result = append(result, eventResponse)
	}

	return result, nil
}

func (s *eventService) FindByID(ctx context.Context, id string, userRole string) (dto.EventResponse, error) {
	event, err := s.eventRepo.GetByID(id)
	if err != nil {
		return dto.EventResponse{}, err
	}

	result := dto.EventResponse{
		ID:          event.ID.String(),
		Name:        event.Name,
		Description: event.Description,
		Price:       event.Price,
		StartDate:   event.StartDate,
		EndDate:     event.EndDate,
	}

	if userRole == constants.ENUM_ROLE_ADMIN {
		if event.Registers == 0 {
			event.Registers = 1
		}

		result.Capacity = event.Capacity
		result.Registers = event.Registers
	}

	return result, nil
}

func (s *eventService) GetPE2Detail(ctx context.Context, userRole string) (dto.EventResponse, error) {
	pe2, err := s.eventRepo.GetPE2Detail()
	if err != nil {
		return dto.EventResponse{}, err
	}

	result := dto.EventResponse{
		ID:          pe2.ID.String(),
		Name:        pe2.Name,
		Description: pe2.Description,
		Price:       pe2.Price,
		StartDate:   pe2.StartDate,
		EndDate:     pe2.EndDate,
	}

	if userRole == constants.ENUM_ROLE_ADMIN {
		if pe2.Registers == 0 {
			pe2.Registers = 1
		}

		result.Capacity = pe2.Capacity
		result.Registers = pe2.Registers
	}

	return result, nil
}
