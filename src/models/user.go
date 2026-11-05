package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name         string    `gorm:"type:text;not null" json:"name"`
	Username     string    `gorm:"type:text;uniqueIndex;not null" json:"username"`
	Email        string    `gorm:"type:text;uniqueIndex;not null" json:"email"`
	Bio          string    `gorm:"type:text" json:"bio"`
	AvatarKey    string    `gorm:"type:text" json:"avatarKey"`
	PasswordHash string    `gorm:"column:password_hash;type:text;not null" json:"-"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
