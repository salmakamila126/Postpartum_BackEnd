package repository

import (
	"Postpartum_BackEnd/internal/entity"

	"gorm.io/gorm"
)

type IRefreshTokenRepository interface {
	Create(token *entity.RefreshToken) error
	Find(token string) (*entity.RefreshToken, error)
	Delete(token string) error
}

type refreshTokenRepository struct {
	DB *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) IRefreshTokenRepository {
	return &refreshTokenRepository{DB: db}
}

func (r *refreshTokenRepository) Create(token *entity.RefreshToken) error {
	return r.DB.Create(token).Error
}

func (r *refreshTokenRepository) Find(token string) (*entity.RefreshToken, error) {
	var rt entity.RefreshToken
	err := r.DB.Where("token = ?", token).First(&rt).Error
	return &rt, err
}

func (r *refreshTokenRepository) Delete(token string) error {
	return r.DB.Where("token = ?", token).Delete(&entity.RefreshToken{}).Error
}
