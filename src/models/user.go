package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Email     string    `gorm:"type:text;uniqueIndex;not null" json:"email"`
	Name      *string   `gorm:"type:text" json:"name,omitempty"`
	Username  *string   `gorm:"type:text;uniqueIndex" json:"username,omitempty"`
	Bio       *string   `gorm:"type:text" json:"bio,omitempty"`
	AvatarKey *string   `gorm:"type:text" json:"avatarKey,omitempty"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
