package repository

import (
	"github.com/TEDxITS/website-backend-2024/entity"
	"gorm.io/gorm"
)

type (
	RoleRepository interface {
		GetRolebyId(roleId string) (entity.Role, error)
	}

	roleRepository struct {
		db *gorm.DB
	}
)

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{
		db: db,
	}
}

func (r *roleRepository) GetRolebyId(roleId string) (entity.Role, error) {
	var role entity.Role
	if err := r.db.Where("id = ?", roleId).Take(&role).Error; err != nil {
		return entity.Role{}, err
	}
	return role, nil
}
