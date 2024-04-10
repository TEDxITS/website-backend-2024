package seeders

import (
	"reflect"
	"time"

	"github.com/TEDxITS/website-backend-2024/constants"
	"github.com/TEDxITS/website-backend-2024/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func EventSeeder(db *gorm.DB) error {
	hasTable := db.Migrator().HasTable(&entity.Event{})
	if !hasTable {
		if err := db.Migrator().CreateTable(&entity.Event{}); err != nil {
			return err
		}
	}

	True := true
	False := false

	var eventList []entity.Event
	eventList = append(eventList,
		entity.Event{
			ID:        uuid.MustParse("7de24efe-0aec-469a-bf0c-8fa8cae3ff3f"),
			Name:      constants.PE2Name,
			Price:     0,
			WithKit:   &False,
			Capacity:  40,
			Registers: 0,
			EventDate: time.Date(2024, time.April, 24, 12, 0, 0, 0, time.UTC),
			StartDate: time.Date(2024, time.April, 10, 19, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, time.April, 19, 19, 0, 0, 0, time.UTC),
		}, entity.Event{
			ID:        uuid.MustParse("94f90ecf-882a-479b-9b1c-98fc6ed6183b"),
			Name:      constants.MainEventEarlyBirdNoMerch,
			Price:     85000,
			WithKit:   &False,
			Capacity:  15,
			Registers: 0,
			StartDate: time.Date(2024, time.April, 24, 12, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, time.April, 30, 12, 0, 0, 0, time.UTC),
		}, entity.Event{
			ID:        uuid.MustParse("66257cc5-c64a-494b-985a-02a40348ea91"),
			Name:      constants.MainEventPreSaleNoMerch,
			Price:     100000,
			WithKit:   &False,
			Capacity:  35,
			Registers: 0,
			StartDate: time.Date(2024, time.May, 2, 12, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, time.May, 6, 12, 0, 0, 0, time.UTC),
		}, entity.Event{
			ID:        uuid.MustParse("19eee29a-0948-4827-b41a-4b3015b96508"),
			Name:      constants.MainEventNormalNoMerch,
			Price:     135000,
			WithKit:   &False,
			Capacity:  45,
			Registers: 0,
			StartDate: time.Date(2024, time.May, 7, 12, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, time.May, 24, 12, 0, 0, 0, time.UTC),
		}, entity.Event{
			ID:        uuid.MustParse("ff0972d2-250d-49e9-977f-dd11cbfdca8d"),
			Name:      constants.MainEventEarlyBirdWithMerch,
			Price:     105000,
			WithKit:   &True,
			Capacity:  5,
			Registers: 0,
			StartDate: time.Date(2024, time.April, 24, 12, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, time.April, 30, 12, 0, 0, 0, time.UTC),
		}, entity.Event{
			ID:        uuid.MustParse("edef91cb-f43f-4a30-a5f0-43acd0d6853f"),
			Name:      constants.MainEventPreSaleWithMerch,
			Price:     120000,
			WithKit:   &True,
			Capacity:  15,
			Registers: 0,
			StartDate: time.Date(2024, time.May, 2, 12, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, time.May, 6, 12, 0, 0, 0, time.UTC),
		}, entity.Event{
			ID:        uuid.MustParse("0594226b-c314-42e7-9743-c76fdd6b7099"),
			Name:      constants.MainEventNormalWithMerch,
			Price:     155000,
			WithKit:   &True,
			Capacity:  25,
			Registers: 0,
			StartDate: time.Date(2024, time.May, 7, 12, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, time.May, 24, 12, 0, 0, 0, time.UTC),
		})

	for _, data := range eventList {
		event := entity.Event{}
		if err := db.Where("id = ?", data.ID).Take(&event).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
		}

		// persist registers count, able to change other such as dates and capacity
		if !reflect.DeepEqual(event, entity.Event{}) {
			data.Registers = event.Registers
		}

		if err := db.Save(&data).Error; err != nil {
			if err := db.Create(&data).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
