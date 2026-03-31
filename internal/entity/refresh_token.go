package entity

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	Token     string    `gorm:"primaryKey"`
	UserID    uuid.UUID `gorm:"type:char(36);not null"`
	ExpiresAt time.Time
}
