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

	var eventList []entity.Event
	eventList = append(eventList, entity.Event{
		ID:          uuid.MustParse("7de24efe-0aec-469a-bf0c-8fa8cae3ff3f"),
		Name:        constants.PE2Name,
		Description: "This is a Description",
		Price:       2000,
		Capacity:    40,
		Registers:   0,
		StartDate:   time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2024, time.April, 17, 12, 0, 0, 0, time.UTC),
	})

	for _, data := range eventList {
		event := entity.Event{}
		if err := db.Where("id = ?", data.ID).Take(&event).Error; err != nil {
			return err
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
