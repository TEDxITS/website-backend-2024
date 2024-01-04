package seeder

import (
	"github.com/TEDxITS/website-backend-2024/config"
	"github.com/TEDxITS/website-backend-2024/migrations/seeder/seeders"
	"gorm.io/gorm"
)

func main() {
	db := config.SetUpDatabaseConnection()
	RunSeeders(db)
}

func RunSeeders(db *gorm.DB) error {
	if err := seeders.RoleSeeder(db); err != nil {
		return err
	}

	if err := seeders.UserSeeder(db); err != nil {
		return err
	}

	return nil
}
