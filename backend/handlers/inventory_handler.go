package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/toko-berkah/pos-backend/dto"
	"github.com/toko-berkah/pos-backend/services"
	"github.com/toko-berkah/pos-backend/utils"
)

type InventoryHandler struct {
	inventoryService *services.InventoryService
}

func NewInventoryHandler(inventoryService *services.InventoryService) *InventoryHandler {
	return &InventoryHandler{
		inventoryService: inventoryService,
	}
}

// GetInventory handles getting inventory for a product
// @Summary Get product inventory
// @Description Get inventory details for a product
// @Tags Inventory
// @Accept json
// @Produce json
// @Param productID query string true "Product ID"
// @Success 200 {object} utils.Response{data=dto.InventoryResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/inventory [get]
func (h *InventoryHandler) GetInventory(c *gin.Context) {
	productID := c.Query("productID")
	if productID == "" {
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(
			utils.NewValidationError("productID parameter required"),
		))
		return
	}

	result, err := h.inventoryService.GetInventory(productID)
	if err != nil {
		customErr := err.(*utils.CustomError)
		statusCode := http.StatusNotFound
		if customErr.Code == "NOT_FOUND" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, utils.NewErrorResponse(customErr))
		return
	}

	c.JSON(http.StatusOK, utils.NewSuccessResponse(result, "Inventory retrieved successfully"))
}

// GetAllInventory handles getting all inventory with pagination
// @Summary Get all inventory
// @Description Get all inventory with pagination
// @Tags Inventory
// @Accept json
// @Produce json
// @Param limit query int false "Limit (default 20)"
// @Param offset query int false "Offset (default 0)"
// @Success 200 {object} dto.PaginatedResponse{data=[]dto.InventoryResponse}
// @Failure 400 {object} utils.Response
// @Router /api/v1/inventory/all [get]
func (h *InventoryHandler) GetAllInventory(c *gin.Context) {
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

	results, total, err := h.inventoryService.GetAllInventory(limit, offset)
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

// AdjustInventory handles inventory adjustment
// @Summary Adjust inventory
// @Description Adjust inventory (add/remove stock)
// @Tags Inventory
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.AdjustInventoryRequest true "Adjust inventory request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/inventory/adjustment [post]
func (h *InventoryHandler) AdjustInventory(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.NewErrorResponse(utils.NewUnauthorizedError()))
		return
	}

	var req dto.AdjustInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(
			utils.NewValidationError(err.Error()),
		))
		return
	}

	err := h.inventoryService.AdjustInventory(req, userID.(string))
	if err != nil {
		customErr := err.(*utils.CustomError)
		statusCode := http.StatusInternalServerError
		if customErr.Code == "NOT_FOUND" {
			statusCode = http.StatusNotFound
		} else if customErr.Code == "VALIDATION_ERROR" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, utils.NewErrorResponse(customErr))
		return
	}

	c.JSON(http.StatusOK, utils.NewSuccessResponse(nil, "Inventory adjusted successfully"))
}

// GetLowStockProducts handles getting low stock products
// @Summary Get low stock products
// @Description Get products below threshold
// @Tags Inventory
// @Accept json
// @Produce json
// @Security Bearer
// @Param limit query int false "Limit (default 20)"
// @Param offset query int false "Offset (default 0)"
// @Success 200 {object} dto.PaginatedResponse{data=[]dto.InventoryResponse}
// @Failure 400 {object} utils.Response
// @Router /api/v1/inventory/low-stock [get]
func (h *InventoryHandler) GetLowStockProducts(c *gin.Context) {
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

	results, total, err := h.inventoryService.GetLowStockProducts(limit, offset)
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

// GetInventoryLogs handles getting inventory logs
// @Summary Get inventory logs
// @Description Get inventory transaction logs with pagination
// @Tags Inventory
// @Accept json
// @Produce json
// @Security Bearer
// @Param limit query int false "Limit (default 20)"
// @Param offset query int false "Offset (default 0)"
// @Success 200 {object} dto.PaginatedResponse{data=[]models.InventoryLog}
// @Failure 400 {object} utils.Response
// @Router /api/v1/inventory/logs [get]
func (h *InventoryHandler) GetInventoryLogs(c *gin.Context) {
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

	results, total, err := h.inventoryService.GetInventoryLogs(limit, offset)
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

// GetInventoryLogsByProduct handles getting logs for specific product
// @Summary Get product inventory logs
// @Description Get inventory logs for a specific product
// @Tags Inventory
// @Accept json
// @Produce json
// @Security Bearer
// @Param productID query string true "Product ID"
// @Param limit query int false "Limit (default 20)"
// @Param offset query int false "Offset (default 0)"
// @Success 200 {object} dto.PaginatedResponse{data=[]models.InventoryLog}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/inventory/logs/product [get]
func (h *InventoryHandler) GetInventoryLogsByProduct(c *gin.Context) {
	productID := c.Query("productID")
	if productID == "" {
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(
			utils.NewValidationError("productID parameter required"),
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

	results, total, err := h.inventoryService.GetInventoryLogsByProduct(productID, limit, offset)
	if err != nil {
		customErr := err.(*utils.CustomError)
		statusCode := http.StatusInternalServerError
		if customErr.Code == "NOT_FOUND" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, utils.NewErrorResponse(customErr))
		return
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Data:  results,
		Total: total,
		Page:  offset/limit + 1,
		Limit: limit,
	})
}
