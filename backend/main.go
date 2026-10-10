package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/toko-berkah/pos-backend/config"
	"github.com/toko-berkah/pos-backend/routes"
)

func init() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
}

func main() {
	// Initialize config
	cfg := config.LoadConfig()

	// Initialize database
	db := config.InitDatabase(cfg)
	if db == nil {
		log.Fatal("Failed to initialize database")
	}

	// Run migrations
	if err := config.MigrateDatabase(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize Redis
	redisClient := config.InitRedis(cfg)
	if redisClient == nil {
		log.Fatal("Failed to initialize Redis")
	}

	// Setup router
	router := routes.SetupRoutes(db, redisClient)

	// Start server
	address := fmt.Sprintf("%s:%s", cfg.APIHost, cfg.APIPort)
	log.Printf("Server starting on %s", address)
	if err := router.Run(address); err != nil {
		log.Fatal(err)
	}
}
