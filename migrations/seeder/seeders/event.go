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
			ID:        uuid.MustParse(constants.PreEvent2ID),
			Name:      constants.PE2Name,
			Price:     0,
			WithKit:   &False,
			Capacity:  110,
			Registers: 0,
			EventDate: time.Date(2024, time.April, 24, 12, 0, 0, 0, time.UTC),
			StartDate: time.Date(2024, time.April, 10, 19, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, time.April, 18, 00, 0, 0, 0, time.UTC),
		}, entity.Event{
			ID:        uuid.MustParse(constants.MainEventEarlyBirdNoMerchID),
			Name:      constants.MainEventEarlyBirdNoMerch,
			Price:     85000,
			WithKit:   &False,
			Capacity:  constants.MainEventEarlyBirdNoMerchCapacity,
			Registers: 0,
			StartDate: time.Date(2023, time.May, 6, 15, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, time.May, 7, 15, 0, 0, 0, time.UTC),
		}, entity.Event{
			ID:        uuid.MustParse(constants.MainEventPreSaleNoMerchID),
			Name:      constants.MainEventPreSaleNoMerch,
			Price:     100000,
			WithKit:   &False,
			Capacity:  constants.MainEventPreSaleNoMerchCapacity,
			Registers: 0,
			StartDate: time.Date(2024, time.May, 8, 15, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, time.May, 10, 15, 0, 0, 0, time.UTC),
		}, entity.Event{
			ID:        uuid.MustParse(constants.MainEventNormalNoMerchID),
			Name:      constants.MainEventNormalNoMerch,
			Price:     135000,
			WithKit:   &False,
			Capacity:  constants.MainEventNormalNoMerchCapacity,
			Registers: 0,
			StartDate: time.Date(2024, time.May, 15, 12, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, time.May, 24, 12, 0, 0, 0, time.UTC),
		}, entity.Event{
			ID:        uuid.MustParse(constants.MainEventEarlyBirdWithMerchID),
			Name:      constants.MainEventEarlyBirdWithMerch,
			Price:     105000,
			WithKit:   &True,
			Capacity:  constants.MainEventEarlyBirdWithMerchCapacity,
			Registers: 0,
			StartDate: time.Date(2023, time.May, 6, 15, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, time.May, 7, 15, 0, 0, 0, time.UTC),
		}, entity.Event{
			ID:        uuid.MustParse(constants.MainEventPreSaleWithMerchID),
			Name:      constants.MainEventPreSaleWithMerch,
			Price:     120000,
			WithKit:   &True,
			Capacity:  constants.MainEventPreSaleWithMerchCapacity,
			Registers: 0,
			StartDate: time.Date(2024, time.May, 8, 15, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, time.May, 10, 15, 0, 0, 0, time.UTC),
		}, entity.Event{
			ID:        uuid.MustParse(constants.MainEventNormalWithMerchID),
			Name:      constants.MainEventNormalWithMerch,
			Price:     155000,
			WithKit:   &True,
			Capacity:  constants.MainEventNormalWithMerchCapacity,
			Registers: 0,
			StartDate: time.Date(2024, time.May, 15, 12, 0, 0, 0, time.UTC),
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
