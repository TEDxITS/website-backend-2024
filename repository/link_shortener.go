package repository

import (
	"math"

	"github.com/TEDxITS/website-backend-2024/entity"
	"gorm.io/gorm"
)

type (
	LinkShortenerRepository interface {
		CreateLinkShortener(linkShorten entity.LinkShortener) (entity.LinkShortener, error)
		GetLinkShortenerByAlias(alias string) (entity.LinkShortener, error)
		CheckAliasExist(alias string) (bool, error)
		GetAllLinkShortenerPagination(search string, limit int, page int) ([]entity.LinkShortener, int64, int64, error)
	}

	linkShortenerRepo struct {
		db *gorm.DB
	}
)

func NewLinkShortenerRepository(db *gorm.DB) LinkShortenerRepository {
	return &linkShortenerRepo{
		db: db,
	}
}

func (r *linkShortenerRepo) CreateLinkShortener(linkShorten entity.LinkShortener) (entity.LinkShortener, error) {
	if err := r.db.Create(&linkShorten).Error; err != nil {
		return entity.LinkShortener{}, err
	}

	return linkShorten, nil
}

func (r *linkShortenerRepo) GetLinkShortenerByAlias(alias string) (entity.LinkShortener, error) {
	var linkShorten entity.LinkShortener
	if err := r.db.Where("alias = ?", alias).Take(&linkShorten).Error; err != nil {
		return entity.LinkShortener{}, err
	}

	return linkShorten, nil
}

func (r *linkShortenerRepo) CheckAliasExist(alias string) (bool, error) {
	var linkShorten entity.LinkShortener
	if err := r.db.Where("alias = ?", alias).Take(&linkShorten).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

func (r *linkShortenerRepo) GetAllLinkShortenerPagination(search string, limit int, page int) ([]entity.LinkShortener, int64, int64, error) {
	var links []entity.LinkShortener
	var count int64

	if search != "" {
		err := r.db.Model(&entity.LinkShortener{}).Where("alias LIKE ? OR link LIKE ?", "%"+search+"%", "%"+search+"%").Count(&count).Error
		if err != nil {
			return nil, 0, 0, err
		}
	} else {
		err := r.db.Model(&entity.LinkShortener{}).Count(&count).Error
		if err != nil {
			return nil, 0, 0, err
		}
	}

	maxPage := int64(math.Ceil(float64(count) / float64(limit)))
	offset := (page - 1) * limit

	err := r.db.Where("alias LIKE ? OR link LIKE ?", "%"+search+"%", "%"+search+"%").Offset(offset).Limit(limit).Find(&links).Error
	if err != nil {
		return nil, 0, 0, err
	}

	return links, maxPage, count, nil
}
