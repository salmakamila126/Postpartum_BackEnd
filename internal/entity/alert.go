package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AlertRule struct {
	ID          uuid.UUID `gorm:"primaryKey;type:char(36)"`
	Code        string    `gorm:"type:varchar(50);not null;uniqueIndex"`
	Level       string    `gorm:"type:varchar(10);not null"`
	Disease     string    `gorm:"type:varchar(150);not null"`
	Description string    `gorm:"type:varchar(255);not null"`
	IsActive    bool      `gorm:"default:true"`
	IsSystem    bool      `gorm:"default:false"`
}

func (a *AlertRule) BeforeCreate(tx *gorm.DB) error {
	a.ID = uuid.New()
	return nil
}

func (AlertRule) TableName() string { return "alert_rules" }
