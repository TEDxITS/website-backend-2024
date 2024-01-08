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

	// only create if it does not exist
	for _, data := range listRole {
		var role entity.Role
		err := db.Where(&entity.Role{Name: data.Name}).First(&role).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		exist := db.Find(&role, "name = ?", data.Name).RowsAffected
		if exist == 0 {
			if err := db.Create(&data).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
