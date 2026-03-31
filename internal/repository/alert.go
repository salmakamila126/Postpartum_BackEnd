package repository

import (
	"Postpartum_BackEnd/internal/entity"

	"gorm.io/gorm"
)

type IAlertRepository interface {
	FindActive() ([]entity.AlertRule, error)
	CountAll() (int64, error)
	SeedAll(rules []entity.AlertRule) error
}

type alertRepository struct {
	DB *gorm.DB
}

func NewAlertRepository(db *gorm.DB) IAlertRepository {
	return &alertRepository{DB: db}
}

func (r *alertRepository) FindActive() ([]entity.AlertRule, error) {
	var rules []entity.AlertRule
	err := r.DB.
		Where("is_active = ?", true).
		Order("code asc").
		Find(&rules).Error
	return rules, err
}

func (r *alertRepository) CountAll() (int64, error) {
	var count int64
	err := r.DB.Model(&entity.AlertRule{}).Count(&count).Error
	return count, err
}

func (r *alertRepository) SeedAll(rules []entity.AlertRule) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		for i := range rules {
			rule := rules[i]
			if err := tx.
				Where("code = ?", rule.Code).
				Assign(map[string]interface{}{
					"level":       rule.Level,
					"disease":     rule.Disease,
					"description": rule.Description,
					"is_active":   rule.IsActive,
					"is_system":   rule.IsSystem,
				}).
				FirstOrCreate(&rule).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
