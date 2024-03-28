package seeders

import (
	"encoding/json"
	"io"
	"os"

	"github.com/TEDxITS/website-backend-2024/entity"
	"gorm.io/gorm"
)

func EventSeeder(db *gorm.DB) error {
	hasTable := db.Migrator().HasTable(&entity.Event{})
	if !hasTable {
		if err := db.Migrator().CreateTable(&entity.Event{}); err != nil {
			return err
		}
	}

	jsonFile, err := os.Open("./migrations/seeder/json/event.json")
	if err != nil {
		return err
	}
	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	var eventList []entity.Event
	if err := json.Unmarshal(jsonData, &eventList); err != nil {
		return err
	}

	for _, data := range eventList {
		if err := db.Save(&data).Error; err != nil {
			if err := db.Create(&data).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
