package repository

import (
	"Postpartum_BackEnd/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IBabyRepository interface {
	Create(baby *entity.Baby) error
	FindByUserID(userID uuid.UUID) (*entity.Baby, error)
	Update(baby *entity.Baby) error
}

type babyRepository struct {
	DB *gorm.DB
}

func NewBabyRepository(db *gorm.DB) IBabyRepository {
	return &babyRepository{DB: db}
}

func (r *babyRepository) Create(baby *entity.Baby) error {
	return r.DB.Create(baby).Error
}

func (r *babyRepository) FindByUserID(userID uuid.UUID) (*entity.Baby, error) {
	var baby entity.Baby
	err := r.DB.Where("user_id = ?", userID).First(&baby).Error
	return &baby, err
}

func (r *babyRepository) Update(baby *entity.Baby) error {
	return r.DB.Save(baby).Error
}
