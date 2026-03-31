package repository

import (
	"Postpartum_BackEnd/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IUserRepository interface {
	FindByEmail(email string) (*entity.User, error)
	Create(user *entity.User) error
	FindByID(userID uuid.UUID) (*entity.User, error)
	Update(user *entity.User) error
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.DB.Preload("Baby").Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepository) Create(user *entity.User) error {
	return r.DB.Create(user).Error
}

func (r *userRepository) FindByID(userID uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.DB.Preload("Baby").
		Where("user_id = ?", userID).
		First(&user).Error
	return &user, err
}

func (r *userRepository) Update(user *entity.User) error {
	return r.DB.Model(&entity.User{}).
		Where("user_id = ?", user.UserID).
		Updates(map[string]interface{}{
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		}).Error
}
