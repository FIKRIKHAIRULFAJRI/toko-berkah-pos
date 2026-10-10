package services

import (
	"testing"
	"github.com/toko-berkah/pos-backend/dto"
	"github.com/toko-berkah/pos-backend/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type MockProductRepo struct {}
func (m *MockProductRepo) GetProductByID(id string) (*models.Product, error) {
    return &models.Product{ID: id, SKU: "SKU001", Name: "Test Product", IsActive: true, PriceRetail: 10000, PriceWholesale: 8000}, nil
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open sqlite: %v", err)
	}
	db.AutoMigrate(&models.Transaction{}, &models.TransactionItem{}, &models.TransactionPayment{}, &models.Inventory{}, &models.InventoryLog{}, &models.Product{})
	return db
}

func TestCashTransaction_Success(t *testing.T) {
	db := setupTestDB(t)
    db.Create(&models.Inventory{ProductID: "p1", QtyOnHand: 10})
    db.Create(&models.Product{ID: "p1", SKU: "SKU001", Name: "Test Product", IsActive: true})

	svc := NewTransactionService(db, &MockProductRepo{})

    req := dto.CreateTransactionRequest{
        Items: []dto.TransactionItemRequest{{ProductID: "p1", Qty: 2}},
        PaymentMethod: "CASH",
        AmountPaid: 25000,
    }
    res, err := svc.CreateTransaction(req, "u1")
    if err != nil {
        t.Fatalf("Cash transaction failed: %v", err)
    }

    if res.PaymentStatus != "PAID" {
        t.Errorf("Expected PAID status, got %s", res.PaymentStatus)
    }
    if res.Change != 5000 {
        t.Errorf("Expected change 5000, got %.2f", res.Change)
    }

    var inv models.Inventory
    db.Where("product_id = ?", "p1").First(&inv)
    if inv.QtyOnHand != 8 {
        t.Errorf("Expected stock 8, got %d", inv.QtyOnHand)
    }
}

func TestTransferTransaction_Flow(t *testing.T) {
	db := setupTestDB(t)
    db.Create(&models.Inventory{ProductID: "p1", QtyOnHand: 10})
    db.Create(&models.Product{ID: "p1", SKU: "SKU001", Name: "Test Product", IsActive: true})

	svc := NewTransactionService(db, &MockProductRepo{})

    // 1. Create Transfer (Should be PENDING, stock unchanged)
    req := dto.CreateTransactionRequest{
        Items: []dto.TransactionItemRequest{{ProductID: "p1", Qty: 3}},
        PaymentMethod: "BANK_TRANSFER",
        AmountPaid: 0,
    }
    res, err := svc.CreateTransaction(req, "u1")
    if err != nil {
        t.Fatalf("Transfer creation failed: %v", err)
    }

    if res.PaymentStatus != "PENDING" {
        t.Errorf("Expected PENDING status, got %s", res.PaymentStatus)
    }

    var inv models.Inventory
    db.Where("product_id = ?", "p1").First(&inv)
    if inv.QtyOnHand != 10 {
        t.Errorf("Expected stock to remain 10 for PENDING transfer, got %d", inv.QtyOnHand)
    }

    // 2. Confirm Payment (Should change to PAID, stock reduced)
    confirmedRes, err := svc.ConfirmPayment(res.ID, "u1")
    if err != nil {
        t.Fatalf("Confirm payment failed: %v", err)
    }

    if confirmedRes.PaymentStatus != "PAID" {
        t.Errorf("Expected PAID after confirmation, got %s", confirmedRes.PaymentStatus)
    }

    db.Where("product_id = ?", "p1").First(&inv)
    if inv.QtyOnHand != 7 {
        t.Errorf("Expected stock reduced to 7, got %d", inv.QtyOnHand)
    }

    // 3. Double confirmation should fail
    _, err = svc.ConfirmPayment(res.ID, "u1")
    if err == nil {
        t.Error("Double confirmation should be rejected")
    }
}

func TestCancelTransaction(t *testing.T) {
	db := setupTestDB(t)
    db.Create(&models.Inventory{ProductID: "p1", QtyOnHand: 10})
    db.Create(&models.Product{ID: "p1", SKU: "SKU001", Name: "Test Product", IsActive: true})

	svc := NewTransactionService(db, &MockProductRepo{})

    req := dto.CreateTransactionRequest{
        Items: []dto.TransactionItemRequest{{ProductID: "p1", Qty: 2}},
        PaymentMethod: "BANK_TRANSFER",
    }
    res, err := svc.CreateTransaction(req, "u1")
    if err != nil {
        t.Fatalf("Failed to create tx: %v", err)
    }

    // Cancel PENDING
    cancelled, err := svc.CancelTransaction(res.ID, "u1")
    if err != nil {
        t.Fatalf("Cancel failed: %v", err)
    }
    if cancelled.PaymentStatus != "CANCELLED" {
        t.Errorf("Expected CANCELLED, got %s", cancelled.PaymentStatus)
    }

    // Trying to cancel again or confirm should fail
    _, err = svc.ConfirmPayment(res.ID, "u1")
    if err == nil {
        t.Error("Should not be able to confirm a cancelled transaction")
    }
}
