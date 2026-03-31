package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Psychologist struct {
	ID           uuid.UUID `gorm:"primaryKey;type:char(36)"`
	Name         string    `gorm:"type:varchar(150);not null"`
	Title        string    `gorm:"type:varchar(100);not null"`
	Job          string    `gorm:"type:varchar(100);not null"`
	ExperienceYr int       `gorm:"not null"`
	PriceIDR     int       `gorm:"column:price_idr;not null"`
	PhotoURL     string    `gorm:"type:varchar(255)"`

	Schedules []PsychologistSchedule `gorm:"foreignKey:PsychologistID;constraint:OnDelete:CASCADE"`
}

func (p *Psychologist) BeforeCreate(tx *gorm.DB) error {
	p.ID = uuid.New()
	return nil
}

func (Psychologist) TableName() string { return "psychologists" }

type PsychologistSchedule struct {
	ID             uuid.UUID `gorm:"primaryKey;type:char(36)"`
	PsychologistID uuid.UUID `gorm:"type:char(36);not null;index"`
	DayOfWeek      string    `gorm:"type:varchar(20);not null"`
	StartTime      string    `gorm:"type:varchar(10);not null"`
	EndTime        string    `gorm:"type:varchar(10);not null"`
}

func (p *PsychologistSchedule) BeforeCreate(tx *gorm.DB) error {
	p.ID = uuid.New()
	return nil
}

func (PsychologistSchedule) TableName() string { return "psychologist_schedules" }
