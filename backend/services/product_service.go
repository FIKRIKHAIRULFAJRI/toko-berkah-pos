package services

import (
	"fmt"

	"github.com/toko-berkah/pos-backend/dto"
	"github.com/toko-berkah/pos-backend/models"
	"github.com/toko-berkah/pos-backend/repositories"
	"github.com/toko-berkah/pos-backend/utils"
)

type ProductService struct {
	productRepo   *repositories.ProductRepository
	inventoryRepo *repositories.InventoryRepository
}

func NewProductService(productRepo *repositories.ProductRepository, inventoryRepo *repositories.InventoryRepository) *ProductService {
	return &ProductService{
		productRepo:   productRepo,
		inventoryRepo: inventoryRepo,
	}
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(req dto.CreateProductRequest) (*dto.ProductResponse, error) {
	// Check if SKU already exists
	_, err := s.productRepo.GetProductBySKU(req.SKU)
	if err == nil {
		return nil, utils.NewConflictError("SKU already exists")
	}

	// Create product
	product := &models.Product{
		SKU:            req.SKU,
		Name:           req.Name,
		CategoryID:     req.CategoryID,
		Unit:           req.Unit,
		PriceRetail:    req.PriceRetail,
		PriceWholesale: req.PriceWholesale,
		Cost:           req.Cost,
		QtyThreshold:   req.QtyThreshold,
		IsActive:       true,
	}

	if err := s.productRepo.CreateProduct(product); err != nil {
		return nil, utils.NewInternalError(fmt.Sprintf("Failed to create product: %v", err))
	}

	// Create inventory record
	inventory := &models.Inventory{
		ProductID:   product.ID,
		QtyOnHand:   0,
		QtyReserved: 0,
	}
	s.inventoryRepo.CreateInventory(inventory)

	return &dto.ProductResponse{
		ID:             product.ID,
		SKU:            product.SKU,
		Name:           product.Name,
		CategoryID:     product.CategoryID,
		Unit:           product.Unit,
		PriceRetail:    product.PriceRetail,
		PriceWholesale: product.PriceWholesale,
		Cost:           product.Cost,
		QtyThreshold:   product.QtyThreshold,
		IsActive:       product.IsActive,
	}, nil
}

// GetProduct retrieves product by ID
func (s *ProductService) GetProduct(id string) (*dto.ProductResponse, error) {
	product, err := s.productRepo.GetProductByID(id)
	if err != nil {
		return nil, utils.NewNotFoundError("Product")
	}

	return &dto.ProductResponse{
		ID:             product.ID,
		SKU:            product.SKU,
		Name:           product.Name,
		CategoryID:     product.CategoryID,
		Unit:           product.Unit,
		PriceRetail:    product.PriceRetail,
		PriceWholesale: product.PriceWholesale,
		Cost:           product.Cost,
		QtyThreshold:   product.QtyThreshold,
		IsActive:       product.IsActive,
	}, nil
}

// GetProductBySKU retrieves product by SKU
func (s *ProductService) GetProductBySKU(sku string) (*dto.ProductResponse, error) {
	product, err := s.productRepo.GetProductBySKU(sku)
	if err != nil {
		return nil, utils.NewNotFoundError("Product")
	}

	return &dto.ProductResponse{
		ID:             product.ID,
		SKU:            product.SKU,
		Name:           product.Name,
		CategoryID:     product.CategoryID,
		Unit:           product.Unit,
		PriceRetail:    product.PriceRetail,
		PriceWholesale: product.PriceWholesale,
		Cost:           product.Cost,
		QtyThreshold:   product.QtyThreshold,
		IsActive:       product.IsActive,
	}, nil
}

// GetAllProducts retrieves all products with pagination
func (s *ProductService) GetAllProducts(limit int, offset int) ([]dto.ProductResponse, int64, error) {
	products, total, err := s.productRepo.GetAllProducts(limit, offset)
	if err != nil {
		return nil, 0, utils.NewInternalError(fmt.Sprintf("Failed to get products: %v", err))
	}

	var responses []dto.ProductResponse
	for _, product := range products {
		responses = append(responses, dto.ProductResponse{
			ID:             product.ID,
			SKU:            product.SKU,
			Name:           product.Name,
			CategoryID:     product.CategoryID,
			Unit:           product.Unit,
			PriceRetail:    product.PriceRetail,
			PriceWholesale: product.PriceWholesale,
			Cost:           product.Cost,
			QtyThreshold:   product.QtyThreshold,
			IsActive:       product.IsActive,
		})
	}

	return responses, total, nil
}

// SearchProducts searches products
func (s *ProductService) SearchProducts(query string, limit int, offset int) ([]dto.ProductResponse, int64, error) {
	products, total, err := s.productRepo.SearchProducts(query, limit, offset)
	if err != nil {
		return nil, 0, utils.NewInternalError(fmt.Sprintf("Failed to search products: %v", err))
	}

	var responses []dto.ProductResponse
	for _, product := range products {
		responses = append(responses, dto.ProductResponse{
			ID:             product.ID,
			SKU:            product.SKU,
			Name:           product.Name,
			CategoryID:     product.CategoryID,
			Unit:           product.Unit,
			PriceRetail:    product.PriceRetail,
			PriceWholesale: product.PriceWholesale,
			Cost:           product.Cost,
			QtyThreshold:   product.QtyThreshold,
			IsActive:       product.IsActive,
		})
	}

	return responses, total, nil
}

// UpdateProduct updates product information
func (s *ProductService) UpdateProduct(id string, req dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	product, err := s.productRepo.GetProductByID(id)
	if err != nil {
		return nil, utils.NewNotFoundError("Product")
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.PriceRetail > 0 {
		updates["price_retail"] = req.PriceRetail
	}
	if req.PriceWholesale > 0 {
		updates["price_wholesale"] = req.PriceWholesale
	}
	if req.Cost > 0 {
		updates["cost"] = req.Cost
	}
	if req.QtyThreshold >= 0 {
		updates["qty_threshold"] = req.QtyThreshold
	}

	if err := s.productRepo.UpdateProduct(id, updates); err != nil {
		return nil, utils.NewInternalError(fmt.Sprintf("Failed to update product: %v", err))
	}

	// Fetch updated product
	product, _ = s.productRepo.GetProductByID(id)

	return &dto.ProductResponse{
		ID:             product.ID,
		SKU:            product.SKU,
		Name:           product.Name,
		CategoryID:     product.CategoryID,
		Unit:           product.Unit,
		PriceRetail:    product.PriceRetail,
		PriceWholesale: product.PriceWholesale,
		Cost:           product.Cost,
		QtyThreshold:   product.QtyThreshold,
		IsActive:       product.IsActive,
	}, nil
}
