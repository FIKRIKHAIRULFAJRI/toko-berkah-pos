// Auth types
export interface User {
  id: string
  username: string
  email: string
  role: 'kasir' | 'owner'
  is_active: boolean
}

export interface LoginResponse {
  access_token: string
  refresh_token: string
  user: User
}

export interface RefreshTokenResponse {
  access_token: string
}

// Product types
export interface Product {
  id: string
  sku: string
  name: string
  category_id: string
  unit: string
  price_retail: number
  price_wholesale: number
  cost: number
  qty_threshold: number
  is_active: boolean
}

export interface CreateProductRequest {
  sku: string
  name: string
  category_id: string
  unit: string
  price_retail: number
  price_wholesale: number
  cost: number
  qty_threshold: number
}

export interface UpdateProductRequest {
  name?: string
  price_retail?: number
  price_wholesale?: number
  cost?: number
  qty_threshold?: number
  is_active?: boolean
}

export interface ProductResponse extends Product {}

// Inventory types
export interface Inventory {
  product_id: string
  product_name: string
  qty_on_hand: number
  qty_reserved: number
  qty_threshold: number
  status: 'low' | 'normal'
}

export interface InventoryResponse extends Inventory {}

export interface AdjustInventoryRequest {
  product_id: string
  qty_change: number
  type: 'adjustment' | 'return' | 'restocking'
  notes?: string
}

// Transaction types
export interface TransactionItem {
  product_id: string
  product_name: string
  qty: number
  unit_price: number
  subtotal: number
}

export interface Transaction {
  id: string
  no_ref: string
  transaction_date: string
  type: 'retail' | 'wholesale'
  total_items: number
  subtotal: number
  discount_amount: number
  total_amount: number
  payment_method: 'cash' | 'transfer'
  payment_status: 'pending' | 'completed' | 'failed'
  items: TransactionItem[]
  created_at: string
}

export interface TransactionResponse extends Transaction {}

export interface CreateTransactionRequest {
  items: Array<{
    product_id: string
    qty: number
  }>
  payment_method: 'cash' | 'transfer'
  amount_paid: number
  reference_number?: string
  notes?: string
}

// Report types
export interface DashboardMetrics {
  total_sales: number
  total_profit: number
  transaction_count: number
  comparison_prev_day: {
    sales_diff_percent: number
    profit_diff_percent: number
  }
}

export interface SalesBreakdown {
  retail: {
    revenue: number
    profit: number
    margin_percent: number
    transaction_count: number
  }
  wholesale: {
    revenue: number
    profit: number
    margin_percent: number
    transaction_count: number
  }
}

export interface BestsellerItem {
  rank: number
  product_name: string
  qty_sold: number
  revenue: number
}

// Pagination
export interface PaginatedResponse {
  data: any[]
  total: number
  page: number
  limit: number
}
