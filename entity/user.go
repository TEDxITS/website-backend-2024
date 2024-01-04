package entity

import (
	"github.com/TEDxITS/website-backend-2024/helpers"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID         uuid.UUID `json:"id" form:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()" `
	Name       string    `json:"name" form:"name"`
	Email      string    `json:"email" form:"email"`
	Password   string    `json:"password" form:"password"`
	IsVerified bool      `json:"is_verified" form:"is_verified"`
	RoleID     string    `json:"role_id" form:"role_id" gorm:"foreignKey" `
	Role       *Role     `json:"role,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" `

	Timestamp
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	var err error
	u.Password, err = helpers.HashPassword(u.Password)
	if err != nil {
		return err
	}
	return nil
}
