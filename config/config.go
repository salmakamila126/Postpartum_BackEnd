package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func NewConfig() {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Printf("failed to load .env file: %v", err)
		}
	}
}
