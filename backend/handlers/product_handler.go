package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/toko-berkah/pos-backend/dto"
	"github.com/toko-berkah/pos-backend/services"
	"github.com/toko-berkah/pos-backend/utils"
)

type ProductHandler struct {
	productService *services.ProductService
}

func NewProductHandler(productService *services.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// CreateProduct handles product creation
// @Summary Create new product
// @Description Create new product (owner only)
// @Tags Products
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.CreateProductRequest true "Create product request"
// @Success 201 {object} utils.Response{data=dto.ProductResponse}
// @Failure 400 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Router /api/v1/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(
			utils.NewValidationError(err.Error()),
		))
		return
	}

	result, err := h.productService.CreateProduct(req)
	if err != nil {
		customErr := err.(*utils.CustomError)
		statusCode := http.StatusInternalServerError
		if customErr.Code == "CONFLICT" {
			statusCode = http.StatusConflict
		} else if customErr.Code == "VALIDATION_ERROR" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, utils.NewErrorResponse(customErr))
		return
	}

	c.JSON(http.StatusCreated, utils.NewSuccessResponse(result, "Product created successfully"))
}

// GetProduct handles getting single product
// @Summary Get product by ID
// @Description Get product details by ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} utils.Response{data=dto.ProductResponse}
// @Failure 404 {object} utils.Response
// @Router /api/v1/products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	productID := c.Param("id")

	result, err := h.productService.GetProduct(productID)
	if err != nil {
		customErr := err.(*utils.CustomError)
		statusCode := http.StatusNotFound
		if customErr.Code == "NOT_FOUND" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, utils.NewErrorResponse(customErr))
		return
	}

	c.JSON(http.StatusOK, utils.NewSuccessResponse(result, "Product retrieved successfully"))
}

// GetProductBySKU handles getting product by SKU
// @Summary Get product by SKU
// @Description Get product details by SKU (for barcode scanning)
// @Tags Products
// @Accept json
// @Produce json
// @Param sku query string true "Product SKU"
// @Success 200 {object} utils.Response{data=dto.ProductResponse}
// @Failure 404 {object} utils.Response
// @Router /api/v1/products/sku [get]
func (h *ProductHandler) GetProductBySKU(c *gin.Context) {
	sku := c.Query("sku")
	if sku == "" {
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(
			utils.NewValidationError("SKU parameter required"),
		))
		return
	}

	result, err := h.productService.GetProductBySKU(sku)
	if err != nil {
		customErr := err.(*utils.CustomError)
		statusCode := http.StatusNotFound
		if customErr.Code == "NOT_FOUND" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, utils.NewErrorResponse(customErr))
		return
	}

	c.JSON(http.StatusOK, utils.NewSuccessResponse(result, "Product retrieved successfully"))
}

// GetAllProducts handles getting all products with pagination
// @Summary Get all products
// @Description Get all products with pagination
// @Tags Products
// @Accept json
// @Produce json
// @Param limit query int false "Limit (default 20)"
// @Param offset query int false "Offset (default 0)"
// @Success 200 {object} dto.PaginatedResponse{data=[]dto.ProductResponse}
// @Failure 400 {object} utils.Response
// @Router /api/v1/products [get]
func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	limit := 20
	offset := 0

	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}

	if o := c.Query("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}

	results, total, err := h.productService.GetAllProducts(limit, offset)
	if err != nil {
		customErr := err.(*utils.CustomError)
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(customErr))
		return
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Data:  results,
		Total: total,
		Page:  offset/limit + 1,
		Limit: limit,
	})
}

// SearchProducts handles product search
// @Summary Search products
// @Description Search products by name or SKU
// @Tags Products
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Limit (default 20)"
// @Param offset query int false "Offset (default 0)"
// @Success 200 {object} dto.PaginatedResponse{data=[]dto.ProductResponse}
// @Failure 400 {object} utils.Response
// @Router /api/v1/products/search [get]
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(
			utils.NewValidationError("Search query required"),
		))
		return
	}

	limit := 20
	offset := 0

	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}

	if o := c.Query("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}

	results, total, err := h.productService.SearchProducts(query, limit, offset)
	if err != nil {
		customErr := err.(*utils.CustomError)
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(customErr))
		return
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Data:  results,
		Total: total,
		Page:  offset/limit + 1,
		Limit: limit,
	})
}

// UpdateProduct handles product update
// @Summary Update product
// @Description Update product information (owner only)
// @Tags Products
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Product ID"
// @Param request body dto.UpdateProductRequest true "Update product request"
// @Success 200 {object} utils.Response{data=dto.ProductResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	productID := c.Param("id")

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(
			utils.NewValidationError(err.Error()),
		))
		return
	}

	result, err := h.productService.UpdateProduct(productID, req)
	if err != nil {
		customErr := err.(*utils.CustomError)
		statusCode := http.StatusInternalServerError
		if customErr.Code == "NOT_FOUND" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, utils.NewErrorResponse(customErr))
		return
	}

	c.JSON(http.StatusOK, utils.NewSuccessResponse(result, "Product updated successfully"))
}
