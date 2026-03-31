package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SymptomLog struct {
	ID         uuid.UUID `gorm:"primaryKey;type:char(36)"`
	UserID     uuid.UUID `gorm:"type:char(36);not null;index"`
	Date       time.Time `gorm:"not null;index"`
	IsBackdate bool      `gorm:"default:false"`

	PhysicalData   string `gorm:"type:text"`
	LastAlertLevel string `gorm:"type:varchar(20)"`
	LastAlertData  string `gorm:"type:longtext"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	Bleedings []BleedingLog `gorm:"foreignKey:LogID;constraint:OnDelete:CASCADE"`
	Moods     []MoodLog     `gorm:"foreignKey:LogID;constraint:OnDelete:CASCADE"`
}

func (s *SymptomLog) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.New()
	return
}

func (SymptomLog) TableName() string {
	return "symptom_logs"
}

type BleedingLog struct {
	ID    uuid.UUID `gorm:"primaryKey;type:char(36)"`
	LogID uuid.UUID `gorm:"type:char(36);not null;index"`

	PadUsage   string `gorm:"type:varchar(20)"`
	ClotSize   string `gorm:"type:varchar(20)"`
	BloodColor string `gorm:"type:varchar(20)"`
	Smell      string `gorm:"type:varchar(20)"`
}

func (b *BleedingLog) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New()
	return
}

func (BleedingLog) TableName() string { return "bleeding_logs" }

type MoodLog struct {
	ID    uuid.UUID `gorm:"primaryKey;type:char(36)"`
	LogID uuid.UUID `gorm:"type:char(36);not null;index"`

	Type string `gorm:"type:varchar(50);not null"`
}

func (m *MoodLog) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}

func (MoodLog) TableName() string {
	return "mood_logs"
}
