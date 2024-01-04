package seeders

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/TEDxITS/website-backend-2024/entity"
	"gorm.io/gorm"
)

func RoleSeeder(db *gorm.DB) error {
	hasTable := db.Migrator().HasTable(&entity.Role{})
	if !hasTable {
		if err := db.Migrator().CreateTable(&entity.Role{}); err != nil {
			return err
		}
	}

	jsonFile, err := os.Open("./migrations/seeder/json/role.json")
	if err != nil {
		return err
	}
	jsonData, _ := io.ReadAll(jsonFile)

	var listRole []entity.Role
	json.Unmarshal(jsonData, &listRole)

	for _, data := range listRole {
		if err := db.Save(&data).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := db.Create(&data).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}
