package models

import (
	"time"

	"github.com/google/uuid"
)

type Friendship struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	RequesterID uuid.UUID `gorm:"type:uuid;not null;index" json:"requesterId"`
	ReceiverID  uuid.UUID `gorm:"type:uuid;not null;index" json:"receiverId"`
	Status      string    `gorm:"type:text;not null;default:'pending'" json:"status"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Requester User `gorm:"foreignKey:RequesterID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Receiver  User `gorm:"foreignKey:ReceiverID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}
