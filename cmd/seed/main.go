package main

import (
	"Postpartum_BackEnd/config"
	"Postpartum_BackEnd/internal/repository"
	"Postpartum_BackEnd/internal/usecase"
	"Postpartum_BackEnd/pkg/logger"
	"Postpartum_BackEnd/pkg/mysql"
	"log"
)

func main() {
	config.NewConfig()
	logger.Init()
	defer logger.Log.Sync()

	db := mysql.StartMySQL()
	repo := repository.NewRepository(db)
	sleepCfg := config.NewSleepConfig()
	uc := usecase.NewUsecase(repo, db, sleepCfg, nil)

	if err := uc.PsychologistUsecase.SeedIfEmpty(); err != nil {
		log.Fatalf("psychologist seed failed: %v", err)
	}

	if err := uc.SymptomUsecase.SeedAlertRulesIfEmpty(); err != nil {
		log.Fatalf("alert rule seed failed: %v", err)
	}

	log.Println("Seed completed successfully")
}
