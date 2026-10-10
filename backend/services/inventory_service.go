package services

import (
	"fmt"
	"time"

	"github.com/toko-berkah/pos-backend/dto"
	"github.com/toko-berkah/pos-backend/models"
	"github.com/toko-berkah/pos-backend/repositories"
	"github.com/toko-berkah/pos-backend/utils"
)

type InventoryService struct {
	inventoryRepo *repositories.InventoryRepository
	productRepo   *repositories.ProductRepository
}

func NewInventoryService(inventoryRepo *repositories.InventoryRepository, productRepo *repositories.ProductRepository) *InventoryService {
	return &InventoryService{
		inventoryRepo: inventoryRepo,
		productRepo:   productRepo,
	}
}

// GetInventory retrieves inventory for a product
func (s *InventoryService) GetInventory(productID string) (*dto.InventoryResponse, error) {
	inventory, err := s.inventoryRepo.GetInventoryByProductID(productID)
	if err != nil {
		return nil, utils.NewNotFoundError("Inventory")
	}

	status := "normal"
	if inventory.QtyOnHand <= inventory.Product.QtyThreshold {
		status = "low"
	}

	return &dto.InventoryResponse{
		ProductID:    inventory.ProductID,
		ProductName:  inventory.Product.Name,
		QtyOnHand:    inventory.QtyOnHand,
		QtyReserved:  inventory.QtyReserved,
		QtyThreshold: inventory.Product.QtyThreshold,
		Status:       status,
	}, nil
}

// GetAllInventory retrieves all inventory with pagination
func (s *InventoryService) GetAllInventory(limit int, offset int) ([]dto.InventoryResponse, int64, error) {
	inventory, total, err := s.inventoryRepo.GetAllInventory(limit, offset)
	if err != nil {
		return nil, 0, utils.NewInternalError(fmt.Sprintf("Failed to get inventory: %v", err))
	}

	var responses []dto.InventoryResponse
	for _, inv := range inventory {
		status := "normal"
		if inv.QtyOnHand <= inv.Product.QtyThreshold {
			status = "low"
		}

		responses = append(responses, dto.InventoryResponse{
			ProductID:    inv.ProductID,
			ProductName:  inv.Product.Name,
			QtyOnHand:    inv.QtyOnHand,
			QtyReserved:  inv.QtyReserved,
			QtyThreshold: inv.Product.QtyThreshold,
			Status:       status,
		})
	}

	return responses, total, nil
}

// AdjustInventory adjusts inventory and creates log
func (s *InventoryService) AdjustInventory(req dto.AdjustInventoryRequest, userID string) error {
	// Get product
	_, err := s.productRepo.GetProductByID(req.ProductID)
	if err != nil {
		return utils.NewNotFoundError("Product")
	}

	// Get current inventory
	inventory, err := s.inventoryRepo.GetInventoryByProductID(req.ProductID)
	if err != nil {
		return utils.NewNotFoundError("Inventory")
	}

	// Store previous qty
	qtyBefore := inventory.QtyOnHand

	// Update inventory
	newQty := qtyBefore + req.QtyChange
	if newQty < 0 {
		return utils.NewValidationError("Insufficient inventory")
	}

	updates := map[string]interface{}{
		"qty_on_hand": newQty,
		"last_updated": time.Now(),
		"updated_by": userID,
	}

	if err := s.inventoryRepo.UpdateInventory(req.ProductID, updates); err != nil {
		return utils.NewInternalError(fmt.Sprintf("Failed to update inventory: %v", err))
	}

	// Create inventory log
	log := &models.InventoryLog{
		ProductID:  req.ProductID,
		Type:       req.Type,
		QtyChange:  req.QtyChange,
		QtyBefore:  qtyBefore,
		QtyAfter:   newQty,
		Notes:      req.Notes,
		UserID:     userID,
	}

	if err := s.inventoryRepo.CreateInventoryLog(log); err != nil {
		// Log the error but don't fail the inventory adjustment
		fmt.Printf("Failed to create inventory log: %v\n", err)
	}

	return nil
}

// GetLowStockProducts retrieves products below threshold
func (s *InventoryService) GetLowStockProducts(limit int, offset int) ([]dto.InventoryResponse, int64, error) {
	inventory, total, err := s.inventoryRepo.GetLowStockProducts(limit, offset)
	if err != nil {
		return nil, 0, utils.NewInternalError(fmt.Sprintf("Failed to get low stock products: %v", err))
	}

	var responses []dto.InventoryResponse
	for _, inv := range inventory {
		responses = append(responses, dto.InventoryResponse{
			ProductID:    inv.ProductID,
			ProductName:  inv.Product.Name,
			QtyOnHand:    inv.QtyOnHand,
			QtyReserved:  inv.QtyReserved,
			QtyThreshold: inv.Product.QtyThreshold,
			Status:       "low",
		})
	}

	return responses, total, nil
}

// GetInventoryLogs retrieves inventory logs with pagination
func (s *InventoryService) GetInventoryLogs(limit int, offset int) ([]models.InventoryLog, int64, error) {
	logs, total, err := s.inventoryRepo.GetInventoryLogs(limit, offset)
	if err != nil {
		return nil, 0, utils.NewInternalError(fmt.Sprintf("Failed to get inventory logs: %v", err))
	}

	return logs, total, nil
}

// GetInventoryLogsByProduct retrieves logs for specific product
func (s *InventoryService) GetInventoryLogsByProduct(productID string, limit int, offset int) ([]models.InventoryLog, int64, error) {
	// Verify product exists
	_, err := s.productRepo.GetProductByID(productID)
	if err != nil {
		return nil, 0, utils.NewNotFoundError("Product")
	}

	logs, total, err := s.inventoryRepo.GetInventoryLogsByProduct(productID, limit, offset)
	if err != nil {
		return nil, 0, utils.NewInternalError(fmt.Sprintf("Failed to get inventory logs: %v", err))
	}

	return logs, total, nil
}
