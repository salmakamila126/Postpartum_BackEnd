package usecase

import (
	"Postpartum_BackEnd/config"
	"Postpartum_BackEnd/internal/repository"
	"Postpartum_BackEnd/pkg/cache"
	"Postpartum_BackEnd/pkg/logger"

	"gorm.io/gorm"
)

type Usecase struct {
	AuthUsecase         *AuthUsecase
	UserUsecase         *UserUsecase
	SleepUsecase        *SleepUsecase
	SymptomUsecase      *SymptomUsecase
	PsychologistUsecase *PsychologistUsecase
}

func NewUsecase(repo *repository.Repository, db *gorm.DB, cfg config.SleepConfig, appCache *cache.Cache) *Usecase {
	return &Usecase{
		AuthUsecase:         NewAuthUsecase(repo, db),
		UserUsecase:         NewUserUsecase(repo),
		SleepUsecase:        NewSleepUsecase(repo, db, cfg, logger.Log),
		SymptomUsecase:      NewSymptomUsecase(repo, appCache),
		PsychologistUsecase: NewPsychologistUsecase(repo, appCache),
	}
}
