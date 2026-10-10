package repositories

import (
	"github.com/toko-berkah/pos-backend/models"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// CreateProduct creates a new product
func (r *ProductRepository) CreateProduct(product *models.Product) error {
	return r.db.Create(product).Error
}

// GetProductByID retrieves product by ID
func (r *ProductRepository) GetProductByID(id string) (*models.Product, error) {
	var product models.Product
	if err := r.db.Preload("Category").First(&product, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

// GetProductBySKU retrieves product by SKU
func (r *ProductRepository) GetProductBySKU(sku string) (*models.Product, error) {
	var product models.Product
	if err := r.db.Preload("Category").First(&product, "sku = ?", sku).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

// GetAllProducts retrieves all active products with pagination
func (r *ProductRepository) GetAllProducts(limit int, offset int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.Where("is_active = ?", true)

	if err := query.Model(&models.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Category").Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// GetProductsByCategory retrieves products by category
func (r *ProductRepository) GetProductsByCategory(categoryID string, limit int, offset int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.Where("category_id = ? AND is_active = ?", categoryID, true)

	if err := query.Model(&models.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Category").Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// UpdateProduct updates product information
func (r *ProductRepository) UpdateProduct(id string, updates map[string]interface{}) error {
	return r.db.Model(&models.Product{}).Where("id = ?", id).Updates(updates).Error
}

// SearchProducts searches products by name or SKU
func (r *ProductRepository) SearchProducts(query string, limit int, offset int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	q := r.db.Where("(name ILIKE ? OR sku ILIKE ?) AND is_active = ?", "%"+query+"%", "%"+query+"%", true)

	if err := q.Model(&models.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := q.Preload("Category").Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
