package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SleepSession struct {
	SleepID    uuid.UUID  `gorm:"column:id;primaryKey;type:char(36)"`
	UserID     uuid.UUID  `gorm:"type:char(36);not null"`
	SleepTime  time.Time  `gorm:"not null"`
	WakeTime   *time.Time `gorm:""`
	CreatedAt  time.Time  `gorm:"autoCreateTime"`
	IsBackdate bool       `gorm:"default:false"`
}

func (s *SleepSession) BeforeCreate(tx *gorm.DB) (err error) {
	s.SleepID = uuid.New()
	return
}

func (SleepSession) TableName() string {
	return "sleep_sessions"
}
