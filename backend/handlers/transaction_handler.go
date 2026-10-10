package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/toko-berkah/pos-backend/dto"
	"github.com/toko-berkah/pos-backend/services"
	"github.com/toko-berkah/pos-backend/utils"
)

type TransactionHandler struct {
	transactionService *services.TransactionService
}

func NewTransactionHandler(transactionService *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

// CreateTransaction handles transaction creation
// @Summary Create new transaction
// @Description Create new transaction (checkout)
// @Tags Transactions
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.CreateTransactionRequest true "Create transaction request"
// @Success 201 {object} utils.Response{data=dto.TransactionResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/transactions [post]
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.NewErrorResponse(utils.NewUnauthorizedError()))
		return
	}

	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(
			utils.NewValidationError(err.Error()),
		))
		return
	}

	result, err := h.transactionService.CreateTransaction(req, userID.(string))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, utils.NewSuccessResponse(result, "Transaction created successfully"))
}

// ConfirmPayment handles bank transfer manual confirmation
// @Summary Confirm bank transfer payment
// @Description Manually confirm a pending bank transfer transaction by operator
// @Tags Transactions
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Transaction ID"
// @Success 200 {object} utils.Response{data=dto.TransactionResponse}
// @Failure 404 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Router /api/v1/transactions/{id}/confirm-payment [post]
func (h *TransactionHandler) ConfirmPayment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.NewErrorResponse(utils.NewUnauthorizedError()))
		return
	}

	transactionID := c.Param("id")

	result, err := h.transactionService.ConfirmPayment(transactionID, userID.(string))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, utils.NewSuccessResponse(result, "Payment confirmed successfully"))
}

// CancelTransaction handles pending transaction cancellation
// @Summary Cancel pending transaction
// @Description Cancel a pending transaction (only allowed for PENDING status)
// @Tags Transactions
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Transaction ID"
// @Success 200 {object} utils.Response{data=dto.TransactionResponse}
// @Failure 404 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Router /api/v1/transactions/{id}/cancel [post]
func (h *TransactionHandler) CancelTransaction(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.NewErrorResponse(utils.NewUnauthorizedError()))
		return
	}

	transactionID := c.Param("id")

	result, err := h.transactionService.CancelTransaction(transactionID, userID.(string))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, utils.NewSuccessResponse(result, "Transaction cancelled successfully"))
}

// GetTransaction handles getting single transaction
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	transactionID := c.Param("id")

	result, err := h.transactionService.GetTransaction(transactionID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, utils.NewSuccessResponse(result, "Transaction retrieved successfully"))
}

// GetAllTransactions handles getting all transactions with pagination
func (h *TransactionHandler) GetAllTransactions(c *gin.Context) {
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

	results, total, err := h.transactionService.GetAllTransactions(limit, offset)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Data:  results,
		Total: total,
		Page:  offset/limit + 1,
		Limit: limit,
	})
}

// GetTransactionsByDateRange handles getting transactions within date range
func (h *TransactionHandler) GetTransactionsByDateRange(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(
			utils.NewValidationError("start_date and end_date parameters required"),
		))
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(
			utils.NewValidationError("Invalid start_date format (use YYYY-MM-DD)"),
		))
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(
			utils.NewValidationError("Invalid end_date format (use YYYY-MM-DD)"),
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

	results, total, err := h.transactionService.GetTransactionsByDateRange(startDate, endDate, limit, offset)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Data:  results,
		Total: total,
		Page:  offset/limit + 1,
		Limit: limit,
	})
}

func handleError(c *gin.Context, err error) {
	customErr, ok := err.(*utils.CustomError)
	if !ok {
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(utils.NewInternalError(err.Error())))
		return
	}

	statusCode := http.StatusInternalServerError
	switch customErr.Code {
	case "NOT_FOUND":
		statusCode = http.StatusNotFound
	case "VALIDATION_ERROR":
		statusCode = http.StatusBadRequest
	case "CONFLICT":
		statusCode = http.StatusConflict
	case "UNAUTHORIZED":
		statusCode = http.StatusUnauthorized
	case "FORBIDDEN":
		statusCode = http.StatusForbidden
	}

	c.JSON(statusCode, utils.NewErrorResponse(customErr))
}
