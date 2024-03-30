package repository

import (
	"math"

	"github.com/TEDxITS/website-backend-2024/entity"
	"gorm.io/gorm"
)

type (
	PE2RSVPRepository interface {
		Create(entity.PE2RSVP) (entity.PE2RSVP, error)
		CheckEmailExist(string) (bool, error)
		GetAllPagination(string, int, int) ([]entity.PE2RSVP, int64, int64, error)
	}

	pe2RSVPRepository struct {
		db *gorm.DB
	}
)

func NewPE2RSVPRepository(db *gorm.DB) PE2RSVPRepository {
	return &pe2RSVPRepository{
		db: db,
	}
}

func (r *pe2RSVPRepository) Create(rsvp entity.PE2RSVP) (entity.PE2RSVP, error) {
	if err := r.db.Create(&rsvp).Error; err != nil {
		return entity.PE2RSVP{}, err
	}

	return rsvp, nil
}

func (r *pe2RSVPRepository) CheckEmailExist(email string) (bool, error) {
	var rsvp entity.PE2RSVP
	if err := r.db.Where("email = ?", email).Take(&rsvp).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (r *pe2RSVPRepository) GetAllPagination(search string, limit, page int) ([]entity.PE2RSVP, int64, int64, error) {
	var rsvps []entity.PE2RSVP
	var count int64

	if search != "" {
		err := r.db.Model(&entity.PE2RSVP{}).Where("name LIKE ?", "%"+search+"%").Count(&count).Error
		if err != nil {
			return nil, 0, 0, err
		}
	} else {
		err := r.db.Model(&entity.PE2RSVP{}).Count(&count).Error
		if err != nil {
			return nil, 0, 0, err
		}
	}

	maxPage := int64(math.Ceil(float64(count) / float64(limit)))
	offset := (page - 1) * limit

	err := r.db.Where("name LIKE ?", "%"+search+"%").Offset(offset).Limit(limit).Find(&rsvps).Error
	if err != nil {
		return nil, 0, 0, err
	}

	return rsvps, maxPage, count, nil
}
