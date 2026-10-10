package config

import (
	"log"

	"github.com/toko-berkah/pos-backend/models"
	"gorm.io/gorm"
)

// MigrateDatabase runs all database migrations
func MigrateDatabase(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// 1. Auto migrate models
	if err := db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Product{},
		&models.Inventory{},
		&models.InventoryLog{},
		&models.Transaction{},
		&models.TransactionItem{},
		&models.TransactionPayment{},
		&models.PricingTier{},
		&models.AuditLog{},
	); err != nil {
		log.Printf("Migration failed: %v", err)
		return err
	}

	// 2. Perform controlled data migration for statuses
	if err := migrateTransactionStatuses(db); err != nil {
		log.Printf("Status migration failed: %v", err)
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}

func migrateTransactionStatuses(db *gorm.DB) error {
	var count int64
	// Check if any legacy statuses exist
	db.Model(&models.Transaction{}).Where("payment_status IN (?)", []string{"pending", "completed", "failed"}).Count(&count)
	if count == 0 {
		log.Println("No legacy statuses found, skipping migration.")
		return nil
	}

	log.Printf("Found %d transactions with legacy statuses. Migrating...", count)

	// Explicit mapping
	mapping := map[string]string{
		"pending":   "PENDING",
		"completed": "PAID",
		"failed":    "CANCELLED",
	}

	for old, new := range mapping {
		if err := db.Model(&models.Transaction{}).Where("payment_status = ?", old).Update("payment_status", new).Error; err != nil {
			return err
		}
	}

	log.Println("Transaction statuses migrated successfully.")
	return nil
}
