package repository

import (
	"math"

	"github.com/TEDxITS/website-backend-2024/entity"
	"gorm.io/gorm"
)

type (
	UserRepository interface {
		RegisterUser(user entity.User) (entity.User, error)
		GetUserById(userId string) (entity.User, error)
		GetUserByEmail(email string) (entity.User, error)
		GetAllUserPagination(search string, limit int, page int) ([]entity.User, int64, int64, error)
		CheckEmailExist(email string) (bool, error)
		UpdateUser(user entity.User) (entity.User, error)
	}

	userRepository struct {
		db *gorm.DB
	}
)

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) RegisterUser(user entity.User) (entity.User, error) {
	var role entity.Role
	if err := r.db.Where("name = ?", "user").First(&role).Error; err != nil {
		return entity.User{}, err
	}

	user.RoleID = role.ID.String()
	if err := r.db.Create(&user).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r *userRepository) GetUserById(userId string) (entity.User, error) {
	var user entity.User
	if err := r.db.Where("id = ?", userId).Take(&user).Error; err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (r *userRepository) GetUserByEmail(email string) (entity.User, error) {
	var user entity.User
	if err := r.db.Where("email = ?", email).Take(&user).Error; err != nil {
		return entity.User{}, err
	}

	var role entity.Role
	if err := r.db.Where("id = ?", user.RoleID).Take(&role).Error; err != nil {
		return entity.User{}, err
	}

	user.Role = &role

	return user, nil
}

func (r *userRepository) GetAllUserPagination(search string, limit int, page int) ([]entity.User, int64, int64, error) {
	var users []entity.User
	var count int64

	if search != "" {
		err := r.db.Model(&entity.User{}).Where("name LIKE ?", "%"+search+"%").Count(&count).Error
		if err != nil {
			return nil, 0, 0, err
		}
	} else {
		err := r.db.Model(&entity.User{}).Count(&count).Error
		if err != nil {
			return nil, 0, 0, err
		}
	}

	maxPage := int64(math.Ceil(float64(count) / float64(limit)))
	offset := (page - 1) * limit

	err := r.db.Where("name LIKE ?", "%"+search+"%").Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, 0, 0, err
	}

	return users, maxPage, count, nil
}

func (r *userRepository) CheckEmailExist(email string) (bool, error) {
	var user entity.User
	if err := r.db.Where("email = ?", email).Take(&user).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (r *userRepository) UpdateUser(user entity.User) (entity.User, error) {
	if err := r.db.Updates(&user).Error; err != nil {
		return entity.User{}, err
	}
	return user, nil
}
