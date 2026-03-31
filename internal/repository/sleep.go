package repository

import (
	"Postpartum_BackEnd/internal/entity"
	"Postpartum_BackEnd/pkg/timeutil"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ISleepRepository interface {
	Create(session *entity.SleepSession) error
	CreateBatch(sessions []entity.SleepSession) error
	Update(session *entity.SleepSession) error
	FindByDate(userID uuid.UUID, date time.Time) ([]entity.SleepSession, error)
	FindActiveSession(userID uuid.UUID) (*entity.SleepSession, error)
	FindRecent(userID uuid.UUID, limit int) ([]entity.SleepSession, error)
	FindHistory(userID uuid.UUID) ([]entity.SleepSession, error)
}

type sleepRepository struct {
	DB *gorm.DB
}

func NewSleepRepository(db *gorm.DB) ISleepRepository {
	return &sleepRepository{DB: db}
}

func (r *sleepRepository) Create(session *entity.SleepSession) error {
	return r.DB.Create(session).Error
}

func (r *sleepRepository) CreateBatch(sessions []entity.SleepSession) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&sessions).Error
	})
}

func (r *sleepRepository) Update(session *entity.SleepSession) error {
	return r.DB.Save(session).Error
}

func (r *sleepRepository) FindByDate(userID uuid.UUID, date time.Time) ([]entity.SleepSession, error) {
	var sessions []entity.SleepSession

	startOfDay := timeutil.StartOfDay(date)
	endOfDay := startOfDay.Add(24 * time.Hour)

	err := r.DB.
		Where("user_id = ?", userID).
		Where("sleep_time >= ? AND sleep_time < ?", startOfDay, endOfDay).
		Order("sleep_time asc").
		Find(&sessions).Error

	return sessions, err
}

func (r *sleepRepository) FindActiveSession(userID uuid.UUID) (*entity.SleepSession, error) {
	var session entity.SleepSession

	err := r.DB.
		Where("user_id = ? AND wake_time IS NULL", userID).
		Order("sleep_time desc").
		First(&session).Error

	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *sleepRepository) FindRecent(userID uuid.UUID, limit int) ([]entity.SleepSession, error) {
	var sessions []entity.SleepSession

	err := r.DB.
		Where("user_id = ?", userID).
		Where("wake_time IS NOT NULL").
		Order("sleep_time desc").
		Limit(limit).
		Find(&sessions).Error

	return sessions, err
}

func (r *sleepRepository) FindHistory(userID uuid.UUID) ([]entity.SleepSession, error) {
	var sessions []entity.SleepSession

	err := r.DB.
		Where("user_id = ?", userID).
		Order("sleep_time desc").
		Find(&sessions).Error

	return sessions, err
}
