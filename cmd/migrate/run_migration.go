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

	// Drop & Recreate tables
	log.Println("⚠️  Dropping all tables...")
	if err := database.DropTables(sqlDB); err != nil {
		log.Fatal("Failed to drop tables:", err)
	}

	log.Println("🔧 Running migrations...")
	if err := database.RunMigrations(sqlDB); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	log.Println("✅ Migration completed successfully!")
}