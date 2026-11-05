package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"`
	TokenHash  string         `gorm:"column:token_hash;type:text;uniqueIndex;not null" json:"-"`
	ClientType string         `gorm:"type:text;not null;default:'web'" json:"clientType"`
	DeviceName string         `gorm:"type:text" json:"deviceName"`
	ExpiresAt  time.Time      `gorm:"not null;index" json:"expiresAt"`
	RevokedAt  *time.Time     `json:"-"`
	CreatedAt  time.Time      `gorm:"not null;autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time      `gorm:"not null;autoUpdateTime" json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	User User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}
