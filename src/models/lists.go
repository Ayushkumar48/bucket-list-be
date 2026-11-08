package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type List struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OwnerID     uuid.UUID `gorm:"type:uuid;not null;index" json:"ownerId"`
	Title       string    `gorm:"type:text;not null" json:"title"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`

	CreatedAt time.Time      `gorm:"not null;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"not null;autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Owner User `gorm:"foreignKey:OwnerID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

type ListItem struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ListID      uuid.UUID `gorm:"type:uuid;not null;index" json:"listId"`
	Title       string    `gorm:"type:text;not null" json:"title"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	Position    int       `gorm:"not null;default:0" json:"position"`

	CreatedAt time.Time      `gorm:"not null;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"not null;autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	List List `gorm:"foreignKey:ListID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

type Challenge struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ListID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_list_challenge" json:"listId"`
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_list_challenge" json:"userId"`

	Status      string         `gorm:"type:text;not null;default:'active'" json:"status"`
	StartedAt   time.Time      `gorm:"not null" json:"startedAt"`
	CompletedAt *time.Time     `json:"completedAt,omitempty"`
	CreatedAt   time.Time      `gorm:"not null;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"not null;autoUpdateTime" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	List List `gorm:"foreignKey:ListID;constraint:OnDelete:CASCADE" json:"-"`
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

type ChallengeItem struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ChallengeID uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex:idx_challenge_list_item" json:"challengeId"`
	ListItemID  uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex:idx_challenge_list_item" json:"listItemId"`
	CompletedAt *time.Time     `json:"completedAt,omitempty"`
	CreatedAt   time.Time      `gorm:"not null;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"not null;autoUpdateTime" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Challenge Challenge `gorm:"foreignKey:ChallengeID;constraint:OnDelete:CASCADE" json:"-"`
	ListItem  ListItem  `gorm:"foreignKey:ListItemID;constraint:OnDelete:CASCADE" json:"-"`
}
