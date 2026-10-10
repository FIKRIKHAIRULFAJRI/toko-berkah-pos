package dto

// Auth DTOs
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginResponse struct {
	AccessToken        string       `json:"access_token"`
	RefreshToken       string       `json:"refresh_token"`
	MustChangePassword bool         `json:"must_change_password"`
	User               UserResponse `json:"user"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required,min=8"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8"`
}

type ChangePasswordResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Message      string `json:"message"`
}

// User DTOs
type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role" binding:"required,oneof=kasir owner"`
	IsActive bool   `json:"is_active"`
}

type UpdateUserRequest struct {
	Username string `json:"username" binding:"min=3,max=50"`
	Email    string `json:"email" binding:"email"`
	Role     string `json:"role" binding:"oneof=kasir owner"`
	IsActive bool   `json:"is_active"`
}

// Product DTOs
type ProductResponse struct {
	ID              string `json:"id"`
	SKU             string `json:"sku"`
	Name            string `json:"name"`
	CategoryID      string `json:"category_id"`
	Unit            string `json:"unit"`
	PriceRetail     float64 `json:"price_retail"`
	PriceWholesale  float64 `json:"price_wholesale"`
	Cost            float64 `json:"cost"`
	QtyThreshold    int     `json:"qty_threshold"`
	IsActive        bool    `json:"is_active"`
}

type CreateProductRequest struct {
	SKU             string  `json:"sku" binding:"required"`
	Name            string  `json:"name" binding:"required"`
	CategoryID      string  `json:"category_id" binding:"required"`
	Unit            string  `json:"unit" binding:"required"`
	PriceRetail     float64 `json:"price_retail" binding:"required,gt=0"`
	PriceWholesale  float64 `json:"price_wholesale" binding:"required,gt=0"`
	Cost            float64 `json:"cost" binding:"required,gt=0"`
	QtyThreshold    int     `json:"qty_threshold" binding:"required,gte=0"`
}

type UpdateProductRequest struct {
	Name            string  `json:"name"`
	PriceRetail     float64 `json:"price_retail" binding:"gt=0"`
	PriceWholesale  float64 `json:"price_wholesale" binding:"gt=0"`
	Cost            float64 `json:"cost" binding:"gt=0"`
	QtyThreshold    int     `json:"qty_threshold" binding:"gte=0"`
	IsActive        bool    `json:"is_active"`
}

// Inventory DTOs
type InventoryResponse struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	QtyOnHand   int    `json:"qty_on_hand"`
	QtyReserved int    `json:"qty_reserved"`
	QtyThreshold int   `json:"qty_threshold"`
	Status      string `json:"status"` // low, normal
}

type AdjustInventoryRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	QtyChange int    `json:"qty_change" binding:"required"`
	Type      string `json:"type" binding:"required,oneof=adjustment return restocking"`
	Notes     string `json:"notes"`
}

// Transaction DTOs
type TransactionItemRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Qty       int    `json:"qty" binding:"required,gt=0"`
}

type CreateTransactionRequest struct {
	Items           []TransactionItemRequest `json:"items" binding:"required,min=1"`
	PaymentMethod   string                   `json:"payment_method" binding:"required,oneof=cash transfer"`
	AmountPaid      float64                  `json:"amount_paid" binding:"required,gt=0"`
	ReferenceNumber string                   `json:"reference_number"`
	Notes           string                   `json:"notes"`
}

type TransactionItemResponse struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Qty         int     `json:"qty"`
	UnitPrice   float64 `json:"unit_price"`
	Subtotal    float64 `json:"subtotal"`
}

type TransactionResponse struct {
	ID              string                     `json:"id"`
	NoRef           string                     `json:"no_ref"`
	TransactionDate string                     `json:"transaction_date"`
	Type            string                     `json:"type"`
	TotalItems      int                        `json:"total_items"`
	Subtotal        float64                    `json:"subtotal"`
	DiscountAmount  float64                    `json:"discount_amount"`
	TotalAmount     float64                    `json:"total_amount"`
	AmountPaid      float64                    `json:"amount_paid"`
	Change          float64                    `json:"change"`
	PaymentMethod   string                     `json:"payment_method"`
	PaymentStatus   string                     `json:"payment_status"`
	Items           []TransactionItemResponse  `json:"items"`
	CreatedAt       string                     `json:"created_at"`
}

// Report DTOs
type DashboardMetrics struct {
	TotalSales    float64 `json:"total_sales"`
	TotalProfit   float64 `json:"total_profit"`
	TransactionCount int  `json:"transaction_count"`
	ComparisonPrevDay struct {
		SalesDiffPercent  float64 `json:"sales_diff_percent"`
		ProfitDiffPercent float64 `json:"profit_diff_percent"`
	} `json:"comparison_prev_day"`
}

type SalesBreakdown struct {
	Retail struct {
		Revenue          float64 `json:"revenue"`
		Profit           float64 `json:"profit"`
		MarginPercent    float64 `json:"margin_percent"`
		TransactionCount int     `json:"transaction_count"`
	} `json:"retail"`
	Wholesale struct {
		Revenue          float64 `json:"revenue"`
		Profit           float64 `json:"profit"`
		MarginPercent    float64 `json:"margin_percent"`
		TransactionCount int     `json:"transaction_count"`
	} `json:"wholesale"`
}

type BestsellerItem struct {
	Rank    int     `json:"rank"`
	ProductName string `json:"product_name"`
	QtySold int     `json:"qty_sold"`
	Revenue float64 `json:"revenue"`
}

// Error Response
type ErrorResponse struct {
	Success   bool        `json:"success"`
	Error     ErrorDetail `json:"error"`
	Timestamp string      `json:"timestamp"`
}

type ErrorDetail struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Success Response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// Pagination
type PaginatedResponse struct {
	Data  interface{} `json:"data"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
}
