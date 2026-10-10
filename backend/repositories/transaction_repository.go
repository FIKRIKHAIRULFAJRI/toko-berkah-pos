package repositories

import (
	"time"

	"github.com/toko-berkah/pos-backend/models"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// CreateTransaction creates a new transaction
func (r *TransactionRepository) CreateTransaction(transaction *models.Transaction) error {
	return r.db.Create(transaction).Error
}

// GetTransactionByID retrieves transaction by ID
func (r *TransactionRepository) GetTransactionByID(id string) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := r.db.Preload("Items").Preload("Items.Product").Preload("Payments").Preload("User").
		First(&transaction, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

// GetTransactionByNoRef retrieves transaction by reference number
func (r *TransactionRepository) GetTransactionByNoRef(noRef string) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := r.db.Preload("Items").Preload("Items.Product").Preload("Payments").
		First(&transaction, "no_ref = ?", noRef).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

// GetAllTransactions retrieves all transactions with pagination
func (r *TransactionRepository) GetAllTransactions(limit int, offset int) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	if err := r.db.Model(&models.Transaction{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Preload("Items").Preload("Items.Product").Preload("Payments").Preload("User").
		Limit(limit).Offset(offset).Order("created_at DESC").Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// GetTransactionsByDateRange retrieves transactions within date range
func (r *TransactionRepository) GetTransactionsByDateRange(startDate, endDate time.Time, limit int, offset int) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	query := r.db.Where("transaction_date BETWEEN ? AND ?", startDate, endDate)

	if err := query.Model(&models.Transaction{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Items").Preload("Items.Product").Preload("Payments").Preload("User").
		Limit(limit).Offset(offset).Order("created_at DESC").Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// GetTransactionsByUser retrieves transactions by user
func (r *TransactionRepository) GetTransactionsByUser(userID string, limit int, offset int) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	query := r.db.Where("user_id = ?", userID)

	if err := query.Model(&models.Transaction{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Items").Preload("Items.Product").Preload("Payments").
		Limit(limit).Offset(offset).Order("created_at DESC").Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// GetTransactionsByType retrieves transactions by type (retail/wholesale)
func (r *TransactionRepository) GetTransactionsByType(txType string, startDate, endDate time.Time) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := r.db.Where("type = ? AND transaction_date BETWEEN ? AND ?", txType, startDate, endDate).
		Preload("Items").Preload("Items.Product").Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

// UpdateTransactionPaymentStatus updates payment status
func (r *TransactionRepository) UpdateTransactionPaymentStatus(id string, status string) error {
	return r.db.Model(&models.Transaction{}).Where("id = ?", id).Update("payment_status", status).Error
}

// CreateTransactionItem creates transaction item
func (r *TransactionRepository) CreateTransactionItem(item *models.TransactionItem) error {
	return r.db.Create(item).Error
}

// CreateTransactionPayment creates transaction payment record
func (r *TransactionRepository) CreateTransactionPayment(payment *models.TransactionPayment) error {
	return r.db.Create(payment).Error
}

// GetDailySalesTotal calculates total sales for a date
func (r *TransactionRepository) GetDailySalesTotal(date time.Time) (float64, error) {
	var total float64
	startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endDate := startDate.AddDate(0, 0, 1)

	if err := r.db.Model(&models.Transaction{}).
		Where("transaction_date >= ? AND transaction_date < ? AND payment_status = ?", startDate, endDate, "completed").
		Select("COALESCE(SUM(total_amount), 0)").Row().Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

// GetDailyProfit calculates profit for a date
func (r *TransactionRepository) GetDailyProfit(date time.Time) (float64, error) {
	var transactions []models.Transaction
	startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endDate := startDate.AddDate(0, 0, 1)

	if err := r.db.Where("transaction_date >= ? AND transaction_date < ? AND payment_status = ?", startDate, endDate, "completed").
		Preload("Items").Preload("Items.Product").Find(&transactions).Error; err != nil {
		return 0, err
	}

	var profit float64
	for _, tx := range transactions {
		for _, item := range tx.Items {
			itemProfit := (item.UnitPrice - item.Product.Cost) * float64(item.Qty)
			profit += itemProfit
		}
	}

	return profit, nil
}
