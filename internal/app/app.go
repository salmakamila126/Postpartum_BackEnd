package app

import (
	"Postpartum_BackEnd/config"
	"Postpartum_BackEnd/internal/controller/rest"
	"Postpartum_BackEnd/internal/repository"
	"Postpartum_BackEnd/internal/usecase"
	"Postpartum_BackEnd/pkg/cache"
	httpserver "Postpartum_BackEnd/pkg/gin"
	"Postpartum_BackEnd/pkg/logger"
	"Postpartum_BackEnd/pkg/mysql"
	"log"
	"os"
)

func Run() {
	config.NewConfig()
	logger.Init()

	db := mysql.StartMySQL()
	app := httpserver.Start()
	appCache, err := cache.NewFromEnv()
	if err != nil {
		log.Printf("redis cache disabled: %v", err)
	}

	repo := repository.NewRepository(db)
	sleepCfg := config.NewSleepConfig()
	uc := usecase.NewUsecase(repo, db, sleepCfg, appCache)

	v1 := rest.NewV1(uc)
	rest.NewRouter(app, v1)

	port := os.Getenv("APP_PORT")
	if port == "" {
		log.Fatal("APP_PORT is required")
	}

	if err := app.Run(":" + port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
