package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID                   string    `gorm:"primaryKey" json:"id"`
	Username             string    `gorm:"uniqueIndex" json:"username"`
	Email                string    `gorm:"uniqueIndex" json:"email"`
	Password             string    `json:"-"`
	Role                 string    `json:"role"` // operator
	IsActive             bool      `json:"is_active"`
	MustChangePassword   bool      `json:"must_change_password"` // Force password change on first login
	PasswordVersion      int       `json:"password_version"`     // For token revocation on password change
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

type Category struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"uniqueIndex" json:"name"`
	Description string  `json:"description"`
	CreatedAt time.Time `json:"created_at"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

type Product struct {
	ID               string    `gorm:"primaryKey" json:"id"`
	SKU              string    `gorm:"uniqueIndex" json:"sku"`
	Name             string    `json:"name"`
	CategoryID       string    `json:"category_id"`
	Category         Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Unit             string    `json:"unit"` // pcs, meter, dll
	PriceRetail      float64   `json:"price_retail"`
	PriceWholesale   float64   `json:"price_wholesale"`
	Cost             float64   `json:"cost"` // Only visible to owner
	QtyThreshold     int       `json:"qty_threshold"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

type Inventory struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	ProductID    string    `json:"product_id"`
	Product      Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	QtyOnHand    int       `json:"qty_on_hand"`
	QtyReserved  int       `json:"qty_reserved"`
	LastUpdated  time.Time `json:"last_updated"`
	UpdatedBy    string    `json:"updated_by"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (i *Inventory) BeforeCreate(tx *gorm.DB) error {
	if i.ID == "" {
		i.ID = uuid.New().String()
	}
	return nil
}

type InventoryLog struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	ProductID string    `json:"product_id"`
	Type      string    `json:"type"` // transaksi, adjustment, return
	QtyChange int       `json:"qty_change"`
	QtyBefore int       `json:"qty_before"`
	QtyAfter  int       `json:"qty_after"`
	Notes     string    `json:"notes"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (l *InventoryLog) BeforeCreate(tx *gorm.DB) error {
	if l.ID == "" {
		l.ID = uuid.New().String()
	}
	return nil
}

type Transaction struct {
	ID              string               `gorm:"primaryKey" json:"id"`
	NoRef           string               `gorm:"uniqueIndex" json:"no_ref"`
	UserID          string               `json:"user_id"`
	User            User                 `gorm:"foreignKey:UserID" json:"user,omitempty"`
	TransactionDate time.Time            `json:"transaction_date"`
	Type            string               `json:"type"` // retail, wholesale
	TotalItems      int                  `json:"total_items"`
	Subtotal        float64              `json:"subtotal"`
	DiscountAmount  float64              `json:"discount_amount"`
	TaxAmount       float64              `json:"tax_amount"`
	TotalAmount     float64              `json:"total_amount"`
	PaymentMethod   string               `json:"payment_method"` // cash, transfer
	PaymentStatus   string               `json:"payment_status"` // PENDING, PAID, CANCELLED
	Notes           string               `json:"notes"`
	Items           []TransactionItem    `gorm:"foreignKey:TransactionID" json:"items,omitempty"`
	Payments        []TransactionPayment `gorm:"foreignKey:TransactionID" json:"payments,omitempty"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}

type TransactionItem struct {
	ID              string    `gorm:"primaryKey" json:"id"`
	TransactionID   string    `json:"transaction_id"`
	ProductID       string    `json:"product_id"`
	Product         Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	ProductName     string    `json:"product_name"` // Snapshot
	ProductSKU      string    `json:"product_sku"`  // Snapshot
	Qty             int       `json:"qty"`
	UnitPrice       float64   `json:"unit_price"`
	Subtotal        float64   `json:"subtotal"`
	DiscountPercent float64   `json:"discount_percent"`
	CreatedAt       time.Time `json:"created_at"`
}

func (ti *TransactionItem) BeforeCreate(tx *gorm.DB) error {
	if ti.ID == "" {
		ti.ID = uuid.New().String()
	}
	return nil
}

type TransactionPayment struct {
	ID                string    `gorm:"primaryKey" json:"id"`
	TransactionID     string    `json:"transaction_id"`
	PaymentMethod     string    `json:"payment_method"` // cash, transfer
	AmountPaid        float64   `json:"amount_paid"`
	ReferenceNumber   string    `json:"reference_number"`
	CreatedAt         time.Time `json:"created_at"`
}

func (tp *TransactionPayment) BeforeCreate(tx *gorm.DB) error {
	if tp.ID == "" {
		tp.ID = uuid.New().String()
	}
	return nil
}

type PricingTier struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	ProductID string    `json:"product_id"`
	QtyMin    int       `json:"qty_min"`
	QtyMax    int       `json:"qty_max"`
	Price     float64   `json:"price"`
	TierName  string    `json:"tier_name"`
	CreatedAt time.Time `json:"created_at"`
}

func (pt *PricingTier) BeforeCreate(tx *gorm.DB) error {
	if pt.ID == "" {
		pt.ID = uuid.New().String()
	}
	return nil
}

type AuditLog struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	UserID       string    `json:"user_id"`
	Action       string    `json:"action"`
	ResourceType string    `json:"resource_type"`
	ResourceID   string    `json:"resource_id"`
	Changes      string    `json:"changes"` // JSON string
	IPAddress    string    `json:"ip_address"`
	CreatedAt    time.Time `json:"created_at"`
}

func (al *AuditLog) BeforeCreate(tx *gorm.DB) error {
	if al.ID == "" {
		al.ID = uuid.New().String()
	}
	return nil
}
