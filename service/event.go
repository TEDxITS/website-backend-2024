package service

import (
	"context"

	"github.com/TEDxITS/website-backend-2024/dto"
	"github.com/TEDxITS/website-backend-2024/repository"
)

type (
	EventService interface {
		FindAll(ctx context.Context, userRole string) (dto.EventResponse, error)
		FindByID(ctx context.Context, id string, userRole string) (dto.EventResponse, error)
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

func (s *eventService) FindAll(ctx context.Context, userRole string) (dto.EventResponse, error) {
	events, err := s.eventRepo.FindAll()
	if err != nil {
		return dto.EventResponse{}, err
	}

	var eventResponses []interface{}
	for _, event := range events {
		var eventResponse interface{}
		if userRole == "admin" {
			eventResponse = dto.AdminEventResponse{
				ID:          event.ID.String(),
				Name:        event.Name,
				Description: event.Description,
				Price:       event.Price,
				Capacity:    event.Capacity,
				Registers:   event.Registers,
				StartDate:   event.StartDate,
				EndDate:     event.EndDate,
			}
		} else {
			eventResponse = dto.UserEventResponse{
				ID:          event.ID.String(),
				Name:        event.Name,
				Description: event.Description,
				Price:       event.Price,
				StartDate:   event.StartDate,
				EndDate:     event.EndDate,
			}
		}
		eventResponses = append(eventResponses, eventResponse)
	}

	return dto.EventResponse{Data: eventResponses}, nil
}

func (s *eventService) FindByID(ctx context.Context, id string, userRole string) (dto.EventResponse, error) {
	event, err := s.eventRepo.FindByID(id)
	if err != nil {
		return dto.EventResponse{}, err
	}

	var eventResponse interface{}

	if userRole == "admin" {
		eventResponse = dto.AdminEventResponse{
			ID:          event.ID.String(),
			Name:        event.Name,
			Description: event.Description,
			Price:       event.Price,
			Capacity:    event.Capacity,
			Registers:   event.Registers,
			StartDate:   event.StartDate,
			EndDate:     event.EndDate,
		}
	} else {
		eventResponse = dto.UserEventResponse{
			ID:          event.ID.String(),
			Name:        event.Name,
			Description: event.Description,
			Price:       event.Price,
			StartDate:   event.StartDate,
			EndDate:     event.EndDate,
		}
	}

	return dto.EventResponse{Data: eventResponse}, nil
}
