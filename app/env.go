package app

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func LoadEnv() {
	if os.Getenv("ENV") == "production" {
		log.Println("Running in production mode!")
		return
	}
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
