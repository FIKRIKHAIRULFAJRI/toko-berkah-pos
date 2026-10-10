package repositories

import (
	"github.com/toko-berkah/pos-backend/models"
	"gorm.io/gorm"
)

type InventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

// GetInventoryByProductID retrieves inventory by product ID
func (r *InventoryRepository) GetInventoryByProductID(productID string) (*models.Inventory, error) {
	var inventory models.Inventory
	if err := r.db.Preload("Product").First(&inventory, "product_id = ?", productID).Error; err != nil {
		return nil, err
	}
	return &inventory, nil
}

// GetAllInventory retrieves all inventory with pagination
func (r *InventoryRepository) GetAllInventory(limit int, offset int) ([]models.Inventory, int64, error) {
	var inventory []models.Inventory
	var total int64

	if err := r.db.Model(&models.Inventory{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Preload("Product").Limit(limit).Offset(offset).Find(&inventory).Error; err != nil {
		return nil, 0, err
	}

	return inventory, total, nil
}

// UpdateInventory updates inventory quantities
func (r *InventoryRepository) UpdateInventory(productID string, updates map[string]interface{}) error {
	return r.db.Model(&models.Inventory{}).Where("product_id = ?", productID).Updates(updates).Error
}

// CreateInventory creates new inventory record
func (r *InventoryRepository) CreateInventory(inventory *models.Inventory) error {
	return r.db.Create(inventory).Error
}

// GetLowStockProducts retrieves products below threshold
func (r *InventoryRepository) GetLowStockProducts(limit int, offset int) ([]models.Inventory, int64, error) {
	var inventory []models.Inventory
	var total int64

	query := r.db.Joins("Product").Where("inventory.qty_on_hand <= product.qty_threshold")

	if err := query.Model(&models.Inventory{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Product").Limit(limit).Offset(offset).Find(&inventory).Error; err != nil {
		return nil, 0, err
	}

	return inventory, total, nil
}

// CreateInventoryLog logs inventory changes
func (r *InventoryRepository) CreateInventoryLog(log *models.InventoryLog) error {
	return r.db.Create(log).Error
}

// GetInventoryLogs retrieves inventory logs with pagination
func (r *InventoryRepository) GetInventoryLogs(limit int, offset int) ([]models.InventoryLog, int64, error) {
	var logs []models.InventoryLog
	var total int64

	if err := r.db.Model(&models.InventoryLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetInventoryLogsByProduct retrieves logs for specific product
func (r *InventoryRepository) GetInventoryLogsByProduct(productID string, limit int, offset int) ([]models.InventoryLog, int64, error) {
	var logs []models.InventoryLog
	var total int64

	query := r.db.Where("product_id = ?", productID)

	if err := query.Model(&models.InventoryLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
