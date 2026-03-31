package repository

import (
	"Postpartum_BackEnd/internal/entity"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IPsychologistRepository interface {
	FindAll() ([]entity.Psychologist, error)
	FindByID(id uuid.UUID) (*entity.Psychologist, error)
	CountAll() (int64, error)
	SeedAll(list []entity.Psychologist) error
	UpdatePhotoURL(id uuid.UUID, photoURL string) error
}

type psychologistRepository struct {
	DB *gorm.DB
}

func NewPsychologistRepository(db *gorm.DB) IPsychologistRepository {
	return &psychologistRepository{DB: db}
}

func (r *psychologistRepository) FindAll() ([]entity.Psychologist, error) {
	var list []entity.Psychologist
	err := r.DB.Order("name asc").Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *psychologistRepository) FindByID(id uuid.UUID) (*entity.Psychologist, error) {
	var p entity.Psychologist
	err := r.DB.
		Preload("Schedules").
		Where("id = ?", id).
		First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *psychologistRepository) CountAll() (int64, error) {
	var count int64
	err := r.DB.Model(&entity.Psychologist{}).Count(&count).Error
	return count, err
}

func (r *psychologistRepository) SeedAll(list []entity.Psychologist) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		for i := range list {
			p := &list[i]
			if err := tx.Omit("Schedules").Create(p).Error; err != nil {
				return err
			}
			for j := range p.Schedules {
				p.Schedules[j].PsychologistID = p.ID
				if err := tx.Create(&p.Schedules[j]).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (r *psychologistRepository) UpdatePhotoURL(id uuid.UUID, photoURL string) error {
	return r.DB.Model(&entity.Psychologist{}).
		Where("id = ?", id).
		Update("photo_url", photoURL).Error
}
