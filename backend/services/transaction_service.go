package services

import (
	"fmt"
	"time"

	"github.com/toko-berkah/pos-backend/dto"
	"github.com/toko-berkah/pos-backend/models"
	"github.com/toko-berkah/pos-backend/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ProductRepositoryInterface defines the contract for product repository
type ProductRepositoryInterface interface {
	GetProductByID(id string) (*models.Product, error)
}

type TransactionService struct {
	db          *gorm.DB
	productRepo ProductRepositoryInterface
}

func NewTransactionService(
	db *gorm.DB,
	productRepo ProductRepositoryInterface,
) *TransactionService {
	return &TransactionService{
		db:          db,
		productRepo: productRepo,
	}
}

// CreateTransaction creates a new transaction (CASH or BANK_TRANSFER) atomically
func (s *TransactionService) CreateTransaction(req dto.CreateTransactionRequest, userID string) (*dto.TransactionResponse, error) {
	if len(req.Items) == 0 {
		return nil, utils.NewValidationError("No items in transaction")
	}

	// Aggregate duplicate Product IDs in request
	itemMap := make(map[string]int)
	for _, item := range req.Items {
		itemMap[item.ProductID] += item.Qty
	}

	var aggregatedItems []dto.TransactionItemRequest
	for prodID, qty := range itemMap {
		aggregatedItems = append(aggregatedItems, dto.TransactionItemRequest{
			ProductID: prodID,
			Qty:       qty,
		})
	}

	var transactionRes *dto.TransactionResponse

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var totalAmount float64
		var totalItems int
		var itemsData []models.TransactionItem

		// 1. Validate and lock inventory / products
		for _, item := range aggregatedItems {
			product, err := s.productRepo.GetProductByID(item.ProductID)
			if err != nil || !product.IsActive {
				return utils.NewNotFoundError(fmt.Sprintf("Product ID %s", item.ProductID))
			}

			// Lock inventory row using SELECT FOR UPDATE
			var inventory models.Inventory
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("product_id = ?", item.ProductID).First(&inventory).Error; err != nil {
				return utils.NewNotFoundError(fmt.Sprintf("Inventory for product %s", product.Name))
			}

			// For CASH, check stock immediately. For BANK_TRANSFER, check current stock availability too.
			if inventory.QtyOnHand < item.Qty {
				return utils.NewConflictError(fmt.Sprintf("Insufficient stock for %s (available: %d, requested: %d)", product.Name, inventory.QtyOnHand, item.Qty))
			}

			// Determine price based on quantity (retail vs wholesale)
			var unitPrice float64
			if item.Qty >= 10 {
				unitPrice = product.PriceWholesale
			} else {
				unitPrice = product.PriceRetail
			}

			subtotal := unitPrice * float64(item.Qty)
			totalAmount += subtotal
			totalItems += item.Qty

			itemsData = append(itemsData, models.TransactionItem{
				ProductID:   item.ProductID,
				ProductName: product.Name,
				ProductSKU:  product.SKU,
				Qty:         item.Qty,
				UnitPrice:   unitPrice,
				Subtotal:    subtotal,
			})
		}

		// Validate CASH payment amount
		if req.PaymentMethod == "CASH" && req.AmountPaid < totalAmount {
			return utils.NewValidationError(fmt.Sprintf("Amount paid (%.2f) is less than total amount (%.2f)", req.AmountPaid, totalAmount))
		}

		// Determine transaction type
		txType := "retail"
		if totalItems >= 10 {
			txType = "wholesale"
		}

		noRef := fmt.Sprintf("TX%d", time.Now().UnixNano())
		paymentStatus := "PENDING"
		if req.PaymentMethod == "CASH" {
			paymentStatus = "PAID"
		}

		transaction := &models.Transaction{
			NoRef:           noRef,
			UserID:          userID,
			TransactionDate: time.Now(),
			Type:            txType,
			TotalItems:      totalItems,
			Subtotal:        totalAmount,
			DiscountAmount:  0,
			TaxAmount:       0,
			TotalAmount:     totalAmount,
			PaymentMethod:   req.PaymentMethod,
			PaymentStatus:   paymentStatus,
			Notes:           req.Notes,
		}

		if err := tx.Create(transaction).Error; err != nil {
			return utils.NewInternalError(fmt.Sprintf("Failed to create transaction: %v", err))
		}

		// Create transaction items
		for i := range itemsData {
			itemsData[i].TransactionID = transaction.ID
			if err := tx.Create(&itemsData[i]).Error; err != nil {
				return utils.NewInternalError(fmt.Sprintf("Failed to create transaction item: %v", err))
			}
		}

		// Create payment record
		payment := &models.TransactionPayment{
			TransactionID:   transaction.ID,
			PaymentMethod:   req.PaymentMethod,
			AmountPaid:      req.AmountPaid,
			ReferenceNumber: req.ReferenceNumber,
		}
		if err := tx.Create(payment).Error; err != nil {
			return utils.NewInternalError(fmt.Sprintf("Failed to create payment record: %v", err))
		}

		// If CASH, deduct inventory immediately
		if req.PaymentMethod == "CASH" {
			for _, item := range aggregatedItems {
				var inventory models.Inventory
				if err := tx.Where("product_id = ?", item.ProductID).First(&inventory).Error; err != nil {
					return err
				}

				qtyBefore := inventory.QtyOnHand
				qtyAfter := qtyBefore - item.Qty

				if err := tx.Model(&models.Inventory{}).Where("product_id = ?", item.ProductID).Update("qty_on_hand", qtyAfter).Error; err != nil {
					return utils.NewInternalError("Failed to update inventory")
				}

				invLog := &models.InventoryLog{
					ProductID: item.ProductID,
					Type:      "transaksi",
					QtyChange: -item.Qty,
					QtyBefore: qtyBefore,
					QtyAfter:  qtyAfter,
					Notes:     fmt.Sprintf("Transaction CASH %s", noRef),
					UserID:    userID,
				}
				if err := tx.Create(invLog).Error; err != nil {
					return utils.NewInternalError("Failed to create inventory log")
				}
			}
		}

		// Load full transaction for response
		var loadedTx models.Transaction
		if err := tx.Preload("Items").Preload("Payments").First(&loadedTx, "id = ?", transaction.ID).Error; err != nil {
			return err
		}

		transactionRes = s.buildTransactionResponse(&loadedTx, req.AmountPaid)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return transactionRes, nil
}

// ConfirmPayment handles manual bank transfer confirmation by operator
func (s *TransactionService) ConfirmPayment(transactionID string, userID string) (*dto.TransactionResponse, error) {
	var transactionRes *dto.TransactionResponse

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Lock transaction row
		var transaction models.Transaction
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Preload("Items").Preload("Payments").First(&transaction, "id = ?", transactionID).Error; err != nil {
			return utils.NewNotFoundError("Transaction")
		}

		// 2. Verify status is PENDING and method is transfer
		if transaction.PaymentStatus != "PENDING" {
			return utils.NewConflictError(fmt.Sprintf("Transaction cannot be confirmed because its status is %s", transaction.PaymentStatus))
		}
		if transaction.PaymentMethod != "transfer" && transaction.PaymentMethod != "BANK_TRANSFER" {
			return utils.NewValidationError("Transaction is not a bank transfer payment")
		}

		// 3. Re-check stock for all items with locking
		for _, item := range transaction.Items {
			var inventory models.Inventory
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("product_id = ?", item.ProductID).First(&inventory).Error; err != nil {
				return utils.NewNotFoundError(fmt.Sprintf("Inventory for product ID %s", item.ProductID))
			}

			if inventory.QtyOnHand < item.Qty {
				return utils.NewConflictError(fmt.Sprintf("Cannot confirm payment: Insufficient stock for product %s (available: %d, required: %d). Transaction remains PENDING.", item.ProductName, inventory.QtyOnHand, item.Qty))
			}
		}

		// 4. Deduct inventory and log
		for _, item := range transaction.Items {
			var inventory models.Inventory
			if err := tx.Where("product_id = ?", item.ProductID).First(&inventory).Error; err != nil {
				return err
			}

			qtyBefore := inventory.QtyOnHand
			qtyAfter := qtyBefore - item.Qty

			if err := tx.Model(&models.Inventory{}).Where("product_id = ?", item.ProductID).Update("qty_on_hand", qtyAfter).Error; err != nil {
				return utils.NewInternalError("Failed to update inventory during confirmation")
			}

			invLog := &models.InventoryLog{
				ProductID: item.ProductID,
				Type:      "transaksi",
				QtyChange: -item.Qty,
				QtyBefore: qtyBefore,
				QtyAfter:  qtyAfter,
				Notes:     fmt.Sprintf("Transfer Confirmed TX %s", transaction.NoRef),
				UserID:    userID,
			}
			if err := tx.Create(invLog).Error; err != nil {
				return utils.NewInternalError("Failed to create inventory log during confirmation")
			}
		}

		// 5. Update status to PAID
		if err := tx.Model(&models.Transaction{}).Where("id = ?", transactionID).Update("payment_status", "PAID").Error; err != nil {
			return utils.NewInternalError("Failed to update transaction status")
		}

		transaction.PaymentStatus = "PAID"
		var amountPaid float64
		if len(transaction.Payments) > 0 {
			amountPaid = transaction.Payments[0].AmountPaid
		}

		transactionRes = s.buildTransactionResponse(&transaction, amountPaid)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return transactionRes, nil
}

// CancelTransaction cancels a PENDING transaction
func (s *TransactionService) CancelTransaction(transactionID string, userID string) (*dto.TransactionResponse, error) {
	var transactionRes *dto.TransactionResponse

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var transaction models.Transaction
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Preload("Items").Preload("Payments").First(&transaction, "id = ?", transactionID).Error; err != nil {
			return utils.NewNotFoundError("Transaction")
		}

		if transaction.PaymentStatus != "PENDING" {
			return utils.NewConflictError(fmt.Sprintf("Only PENDING transactions can be cancelled. Current status: %s", transaction.PaymentStatus))
		}

		if err := tx.Model(&models.Transaction{}).Where("id = ?", transactionID).Update("payment_status", "CANCELLED").Error; err != nil {
			return utils.NewInternalError("Failed to cancel transaction")
		}

		transaction.PaymentStatus = "CANCELLED"
		var amountPaid float64
		if len(transaction.Payments) > 0 {
			amountPaid = transaction.Payments[0].AmountPaid
		}

		transactionRes = s.buildTransactionResponse(&transaction, amountPaid)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return transactionRes, nil
}

func (s *TransactionService) GetTransaction(id string) (*dto.TransactionResponse, error) {
	var transaction models.Transaction
	if err := s.db.Preload("Items").Preload("Items.Product").Preload("Payments").First(&transaction, "id = ?", id).Error; err != nil {
		return nil, utils.NewNotFoundError("Transaction")
	}

	var amountPaid float64
	if len(transaction.Payments) > 0 {
		amountPaid = transaction.Payments[0].AmountPaid
	}

	return s.buildTransactionResponse(&transaction, amountPaid), nil
}

func (s *TransactionService) GetAllTransactions(limit int, offset int) ([]dto.TransactionResponse, int64, error) {
	var transactions []models.Transaction
	var total int64

	if err := s.db.Model(&models.Transaction{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := s.db.Preload("Items").Preload("Items.Product").Preload("Payments").Preload("User").
		Limit(limit).Offset(offset).Order("created_at DESC").Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	var responses []dto.TransactionResponse
	for _, tx := range transactions {
		var amountPaid float64
		if len(tx.Payments) > 0 {
			amountPaid = tx.Payments[0].AmountPaid
		}
		responses = append(responses, *s.buildTransactionResponse(&tx, amountPaid))
	}

	return responses, total, nil
}

func (s *TransactionService) GetTransactionsByDateRange(startDate, endDate time.Time, limit int, offset int) ([]dto.TransactionResponse, int64, error) {
	var transactions []models.Transaction
	var total int64

	query := s.db.Where("transaction_date BETWEEN ? AND ?", startDate, endDate)

	if err := query.Model(&models.Transaction{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Items").Preload("Items.Product").Preload("Payments").Preload("User").
		Limit(limit).Offset(offset).Order("created_at DESC").Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	var responses []dto.TransactionResponse
	for _, tx := range transactions {
		var amountPaid float64
		if len(tx.Payments) > 0 {
			amountPaid = tx.Payments[0].AmountPaid
		}
		responses = append(responses, *s.buildTransactionResponse(&tx, amountPaid))
	}

	return responses, total, nil
}

func (s *TransactionService) buildTransactionResponse(transaction *models.Transaction, amountPaid float64) *dto.TransactionResponse {
	var items []dto.TransactionItemResponse
	for _, item := range transaction.Items {
		items = append(items, dto.TransactionItemResponse{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Qty:         item.Qty,
			UnitPrice:   item.UnitPrice,
			Subtotal:    item.Subtotal,
		})
	}

	change := 0.0
	if amountPaid > transaction.TotalAmount {
		change = amountPaid - transaction.TotalAmount
	}

	return &dto.TransactionResponse{
		ID:              transaction.ID,
		NoRef:           transaction.NoRef,
		TransactionDate: transaction.TransactionDate.Format(time.RFC3339),
		Type:            transaction.Type,
		TotalItems:      transaction.TotalItems,
		Subtotal:        transaction.Subtotal,
		DiscountAmount:  transaction.DiscountAmount,
		TotalAmount:     transaction.TotalAmount,
		AmountPaid:      amountPaid,
		Change:          change,
		PaymentMethod:   transaction.PaymentMethod,
		PaymentStatus:   transaction.PaymentStatus,
		Items:           items,
		CreatedAt:       transaction.CreatedAt.Format(time.RFC3339),
	}
}
