package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var AppConfig Config

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	AppConfig = Config{
		DBUrl: os.Getenv("DATABASE_URL"),
	}
}