# Toko Berkah POS - Backend

Backend untuk sistem kasir Toko Berkah menggunakan Golang + Gin.

## Setup

### Prerequisites
- Go 1.21+
- PostgreSQL 14+
- Redis 7+

### Installation

```bash
cd backend
go mod download
```

### Environment Variables

Copy `.env.example` ke `.env` dan sesuaikan:

```bash
cp .env.example .env
```

### Run Development

```bash
# Dengan Docker Compose
docker-compose up -d

# Run backend
go run main.go
```

API akan tersedia di `http://localhost:8080`

## Project Structure

- `main.go` - Entry point
- `config/` - Configuration, database setup, migrations, seed data
- `models/` - GORM models
- `dto/` - Data transfer objects
- `handlers/` - HTTP handlers
- `repositories/` - Database operations
- `services/` - Business logic
- `middleware/` - HTTP middleware (auth, CORS, error handling)
- `utils/` - Utility functions (JWT, error handling)
- `routes/` - Route definitions

## API Documentation

Swagger docs available at `/swagger/index.html` after running the server.

## Core API Endpoints

### Authentication
- `POST /api/v1/auth/login` - Login user
- `POST /api/v1/auth/refresh` - Refresh access token
- `POST /api/v1/auth/register` - Register new user (owner only)
- `GET /api/v1/auth/profile` - Get user profile (authenticated)
- `POST /api/v1/auth/logout` - Logout user (authenticated)

### Products
- `GET /api/v1/products` - Get all products
- `GET /api/v1/products/:id` - Get product by ID
- `GET /api/v1/products/sku?sku=SKU` - Get product by SKU (barcode scanning)
- `GET /api/v1/products/search?q=query` - Search products
- `POST /api/v1/products` - Create product (owner only)
- `PUT /api/v1/products/:id` - Update product (owner only)

### Inventory
- `GET /api/v1/inventory?productID=ID` - Get inventory for product
- `GET /api/v1/inventory/all` - Get all inventory
- `GET /api/v1/inventory/low-stock` - Get low stock products (owner only)
- `POST /api/v1/inventory/adjustment` - Adjust inventory (owner only)
- `GET /api/v1/inventory/logs` - Get inventory logs (owner only)
- `GET /api/v1/inventory/logs/product?productID=ID` - Get logs for product (owner only)

### Transactions
- `POST /api/v1/transactions` - Create transaction (kasir/owner)
- `GET /api/v1/transactions` - Get all transactions (authenticated)
- `GET /api/v1/transactions/:id` - Get transaction by ID (authenticated)
- `GET /api/v1/transactions/date-range?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD` - Get transactions by date range (owner only)

## Default Credentials (for development)

```
Owner:
- Username: owner
- Password: password123
- Email: owner@toko-berkah.com

Kasir:
- Username: kasir
- Password: password123
- Email: kasir@toko-berkah.com
```

## Features Implemented

### Auth Service
- User login with JWT tokens
- Token refresh mechanism
- User registration (owner only)
- Role-based access control (kasir, owner)
- Password hashing with bcrypt

### Product Service
- CRUD operations for products
- SKU-based product lookup (for barcode scanning)
- Product search by name or SKU
- Retail and wholesale pricing
- Inventory tracking integration

### Inventory Service
- Real-time inventory tracking
- Inventory adjustments with logging
- Low stock notifications
- Inventory transaction history
- Types: adjustment, return, restocking, transaksi

### Transaction Service
- Complete transaction creation and processing
- Automatic price calculation (retail vs wholesale based on quantity)
- Payment method handling (cash, transfer)
- Inventory deduction on transaction
- Transaction history and reporting
- Reference number generation

## Architecture

```
Request → Middleware (Auth, CORS) → Handler → Service → Repository → Database
                                        ↓
                                   Validation
                                   Business Logic
```

## Next Steps

The following features are ready for implementation:

1. **Priority 2 - Frontend Integration**
   - Connect auth API to login/logout flows
   - Build kasir (cashier) view with product search and cart management
   - Implement transaction completion and receipt generation
   - Create owner dashboard with sales reports

2. **Priority 3 - Additional Features**
   - Implement report endpoints (dashboard, sales breakdown, bestseller)
   - Add barcode scanning integration
   - Implement discount and pricing tier logic
   - Build inventory adjustment workflows
   - Add audit logging for transactions
   - Implement error handling and validation across both layers