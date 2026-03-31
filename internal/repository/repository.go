package repository

import "gorm.io/gorm"

type Repository struct {
	AlertRepository        IAlertRepository
	UserRepository         IUserRepository
	BabyRepository         IBabyRepository
	RefreshTokenRepository IRefreshTokenRepository
	SleepRepository        ISleepRepository
	SymptomRepository      ISymptomRepository
	PsychologistRepository IPsychologistRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		AlertRepository:        NewAlertRepository(db),
		UserRepository:         NewUserRepository(db),
		BabyRepository:         NewBabyRepository(db),
		RefreshTokenRepository: NewRefreshTokenRepository(db),
		SleepRepository:        NewSleepRepository(db),
		SymptomRepository:      NewSymptomRepository(db),
		PsychologistRepository: NewPsychologistRepository(db),
	}
}
