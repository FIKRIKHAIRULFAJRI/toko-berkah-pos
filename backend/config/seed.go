package config

import (
	"fmt"
	"log"

	"github.com/toko-berkah/pos-backend/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedDatabase seeds initial data for development
// For production, OPERATOR_USERNAME and OPERATOR_PASSWORD must be set via environment variables
func SeedDatabase(db *gorm.DB) error {
	log.Println("Seeding database...")

	// Check if data already exists
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)
	if userCount > 0 {
		log.Println("Database already seeded, skipping...")
		return nil
	}

	// Load config for bootstrap credentials
	cfg := LoadConfig()

	// Validate bootstrap credentials
	if cfg.OperatorUsername == "" || cfg.OperatorPassword == "" {
		return fmt.Errorf("FATAL: Database is empty but OPERATOR_USERNAME and OPERATOR_PASSWORD are not set. " +
			"Please set these environment variables and restart the application. " +
			"Example: export OPERATOR_USERNAME=operator && export OPERATOR_PASSWORD='secure-password'")
	}

	// Hash operator password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cfg.OperatorPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Failed to hash operator password: %v", err)
	}

	// Create single operator user
	operator := &models.User{
		Username:             cfg.OperatorUsername,
		Email:                "operator@toko-berkah.com",
		Password:             string(hashedPassword),
		Role:                 "operator",
		IsActive:             true,
		MustChangePassword:   true, // Force password change on first login
		PasswordVersion:      1,
	}

	if err := db.Create(operator).Error; err != nil {
		log.Printf("Failed to seed operator user: %v", err)
		return err
	}

	log.Printf("Operator user '%s' created successfully. Password must be changed on first login.", cfg.OperatorUsername)

	// Create categories
	categories := []models.Category{
		{Name: "Elektronik", Description: "Produk elektronik"},
		{Name: "Pakaian", Description: "Produk pakaian"},
		{Name: "Makanan", Description: "Produk makanan"},
		{Name: "Minuman", Description: "Produk minuman"},
	}

	if err := db.Create(categories).Error; err != nil {
		log.Printf("Failed to seed categories: %v", err)
		return err
	}

	// Create sample products
	products := []models.Product{
		{
			SKU:            "ELEC001",
			Name:           "Lampu LED 10W",
			CategoryID:     categories[0].ID,
			Unit:           "pcs",
			PriceRetail:    50000,
			PriceWholesale: 40000,
			Cost:           25000,
			QtyThreshold:   5,
			IsActive:       true,
		},
		{
			SKU:            "CLOTH001",
			Name:           "T-Shirt Pria Size M",
			CategoryID:     categories[1].ID,
			Unit:           "pcs",
			PriceRetail:    75000,
			PriceWholesale: 60000,
			Cost:           40000,
			QtyThreshold:   10,
			IsActive:       true,
		},
		{
			SKU:            "FOOD001",
			Name:           "Beras Premium 5kg",
			CategoryID:     categories[2].ID,
			Unit:           "kg",
			PriceRetail:    65000,
			PriceWholesale: 55000,
			Cost:           40000,
			QtyThreshold:   3,
			IsActive:       true,
		},
		{
			SKU:            "DRINK001",
			Name:           "Minuman Botol 1L",
			CategoryID:     categories[3].ID,
			Unit:           "botol",
			PriceRetail:    12000,
			PriceWholesale: 10000,
			Cost:           6000,
			QtyThreshold:   20,
			IsActive:       true,
		},
	}

	if err := db.Create(products).Error; err != nil {
		log.Printf("Failed to seed products: %v", err)
		return err
	}

	// Create inventory records
	var inventories []models.Inventory
	for _, product := range products {
		inventories = append(inventories, models.Inventory{
			ProductID:  product.ID,
			QtyOnHand:  50,
			QtyReserved: 0,
		})
	}

	if err := db.Create(inventories).Error; err != nil {
		log.Printf("Failed to seed inventory: %v", err)
		return err
	}

	log.Println("Database seeding completed successfully")
	return nil
}
