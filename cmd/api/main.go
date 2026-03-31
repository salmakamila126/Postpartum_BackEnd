package main

import (
	"Postpartum_BackEnd/internal/app"
	"Postpartum_BackEnd/pkg/logger"
)

func main() {
	logger.Init()
	defer logger.Log.Sync()

	app.Run()
}
