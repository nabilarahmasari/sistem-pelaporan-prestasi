package main

import (
	"log"
	"project_uas/config"
	"project_uas/database"
)

func main() {
	// Load config
	config.LoadEnv()

	// Connect database
	database.ConnectDatabase()
	sqlDB, err := database.DB.DB()
	if err != nil {
		log.Fatal("Failed to get database connection:", err)
	}

	// Run seeder
	log.Println("🌱 Running seeder...")
	if err := database.RunSeeders(sqlDB); err != nil {
		log.Fatal("Seeder failed:", err)
	}

	log.Println("✅ Seeder completed successfully!")
}