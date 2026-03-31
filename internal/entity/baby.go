package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Baby struct {
	BabyID    uuid.UUID `gorm:"column:baby_id;primaryKey;type:char(36)" json:"baby_id"`
	UserID    uuid.UUID `gorm:"not null;unique;type:char(36)" json:"user_id"`
	BirthDate string    `gorm:"type:date;not null" json:"birth_date"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (b *Baby) BeforeCreate(tx *gorm.DB) (err error) {
	b.BabyID = uuid.New()
	return
}

func (Baby) TableName() string {
	return "baby"
}
