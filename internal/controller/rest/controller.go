package rest

import "Postpartum_BackEnd/internal/usecase"

type V1 struct {
	Auth         *AuthController
	User         *UserController
	Sleep        *SleepController
	Symptom      *SymptomController
	Psychologist *PsychologistController
}

func NewV1(uc *usecase.Usecase) *V1 {
	return &V1{
		Auth:         NewAuthController(uc.AuthUsecase),
		User:         NewUserController(uc.UserUsecase),
		Sleep:        NewSleepController(uc.SleepUsecase),
		Symptom:      NewSymptomController(uc.SymptomUsecase),
		Psychologist: NewPsychologistController(uc.PsychologistUsecase),
	}
}
