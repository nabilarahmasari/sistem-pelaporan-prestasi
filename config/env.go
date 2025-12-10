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
		DBUrl:      os.Getenv("DATABASE_URL"),
		MongoURL:   getEnv("MONGO_URL", "mongodb://localhost:27017"),
		MongoDB:    getEnv("MONGO_DB", "uas_achievements"),
		Port:       getEnv("PORT", "3000"),
		JWTSecret:  getEnv("JWT_SECRET", "default-secret-key"),
	}

	log.Println("Environment variables loaded successfully")
}

// Helper function untuk get env dengan default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}