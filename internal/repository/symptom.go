package repository

import (
	"Postpartum_BackEnd/internal/entity"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ISymptomRepository interface {
	Create(log *entity.SymptomLog) error
	Update(log *entity.SymptomLog) error
	FindByDate(userID uuid.UUID, date time.Time) (*entity.SymptomLog, error)
	FindHistory(userID uuid.UUID) ([]entity.SymptomLog, error)
	FindMoodHistoryWeekly(userID uuid.UUID, weeks int) ([]entity.SymptomLog, error)
}

type symptomRepository struct {
	DB *gorm.DB
}

func NewSymptomRepository(db *gorm.DB) ISymptomRepository {
	return &symptomRepository{DB: db}
}

func (r *symptomRepository) Create(log *entity.SymptomLog) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(log).Error; err != nil {
			return err
		}

		for i := range log.Bleedings {
			log.Bleedings[i].LogID = log.ID
		}
		if len(log.Bleedings) > 0 {
			if err := tx.Create(&log.Bleedings).Error; err != nil {
				return err
			}
		}

		for i := range log.Moods {
			log.Moods[i].LogID = log.ID
		}
		if len(log.Moods) > 0 {
			if err := tx.Create(&log.Moods).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *symptomRepository) Update(log *entity.SymptomLog) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {

		if err := tx.Model(log).Updates(map[string]interface{}{
			"physical_data":    log.PhysicalData,
			"is_backdate":      log.IsBackdate,
			"last_alert_level": log.LastAlertLevel,
			"last_alert_data":  log.LastAlertData,
			"updated_at":       time.Now(),
		}).Error; err != nil {
			return err
		}

		if err := tx.Where("log_id = ?", log.ID).Delete(&entity.BleedingLog{}).Error; err != nil {
			return err
		}
		if err := tx.Where("log_id = ?", log.ID).Delete(&entity.MoodLog{}).Error; err != nil {
			return err
		}

		for i := range log.Bleedings {
			log.Bleedings[i].LogID = log.ID
			log.Bleedings[i].ID = uuid.New()
		}
		if len(log.Bleedings) > 0 {
			if err := tx.Create(&log.Bleedings).Error; err != nil {
				return err
			}
		}

		for i := range log.Moods {
			log.Moods[i].LogID = log.ID
			log.Moods[i].ID = uuid.New()
		}
		if len(log.Moods) > 0 {
			if err := tx.Create(&log.Moods).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *symptomRepository) FindByDate(userID uuid.UUID, date time.Time) (*entity.SymptomLog, error) {
	startOfDay := date.Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)

	var log entity.SymptomLog

	err := r.DB.
		Where("user_id = ?", userID).
		Where("date >= ? AND date < ?", startOfDay, endOfDay).
		Preload("Bleedings").
		Preload("Moods").
		First(&log).Error

	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *symptomRepository) FindHistory(userID uuid.UUID) ([]entity.SymptomLog, error) {
	var logs []entity.SymptomLog

	err := r.DB.
		Where("user_id = ?", userID).
		Order("date desc").
		Preload("Bleedings").
		Preload("Moods").
		Find(&logs).Error

	return logs, err
}

func (r *symptomRepository) FindMoodHistoryWeekly(userID uuid.UUID, weeks int) ([]entity.SymptomLog, error) {
	since := time.Now().AddDate(0, 0, -(weeks * 7))

	var logs []entity.SymptomLog

	err := r.DB.
		Select("id, user_id, date").
		Where("user_id = ? AND date >= ?", userID, since).
		Order("date asc").
		Preload("Moods").
		Find(&logs).Error

	return logs, err
}
