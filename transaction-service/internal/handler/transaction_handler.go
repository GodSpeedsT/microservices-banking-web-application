package handler

import (
	"net/http"
	"strconv"
	"time"

	"transaction-service/internal/model"
	"transaction-service/internal/service"
	"transaction-service/pkg/security"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionService *service.TransactionService
}

func NewTransactionHandler(transactionService *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

type CreateTransactionRequest struct {
	AccountID   string                `json:"account_id" binding:"required"`
	Amount      float64               `json:"amount" binding:"required,gt=0"`
	Currency    string                `json:"currency" binding:"required"`
	Type        model.TransactionType `json:"type" binding:"required"`
	Description string                `json:"description"`
	Reference   string                `json:"reference"`
}

// CreateTransaction создает новую транзакцию
// @Summary Create a new transaction
// @Description Process a financial transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Param request body CreateTransactionRequest true "Transaction request"
// @Security BearerAuth
// @Success 201 {object} service.ProcessTransactionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /transactions [post]
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "VALIDATION_ERROR",
			Message: err.Error(),
		})
		return
	}

	// Получаем access token из заголовка
	accessToken, err := security.ExtractTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "AUTH_ERROR",
			Message: err.Error(),
		})
		return
	}

	// Получаем user ID из токена (из контекста, установленного middleware)
	claims, exists := c.Get("userClaims")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "AUTH_ERROR",
			Message: "User claims not found",
		})
		return
	}

	jwtClaims, ok := claims.(*security.JWTClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "AUTH_ERROR",
			Message: "Invalid token claims",
		})
		return
	}

	serviceReq := &service.ProcessTransactionRequest{
		UserID:      jwtClaims.UserID,
		AccountID:   req.AccountID,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Type:        req.Type,
		Description: req.Description,
		Reference:   req.Reference,
		AccessToken: accessToken,
	}

	result, err := h.transactionService.ProcessTransaction(c.Request.Context(), serviceReq)
	if err != nil {
		statusCode := http.StatusInternalServerError
		switch err {
		case service.ErrInsufficientBalance:
			statusCode = http.StatusBadRequest
		case service.ErrInvalidAmount:
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, ErrorResponse{
			Error:   "TRANSACTION_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GetTransaction возвращает транзакцию по ID
// @Summary Get transaction by ID
// @Description Get transaction details by transaction ID
// @Tags transactions
// @Produce json
// @Param id path string true "Transaction ID"
// @Security BearerAuth
// @Success 200 {object} model.Transaction
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /transactions/{id} [get]
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	transactionID := c.Param("id")
	if transactionID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "VALIDATION_ERROR",
			Message: "Transaction ID is required",
		})
		return
	}

	accessToken, err := security.ExtractTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "AUTH_ERROR",
			Message: err.Error(),
		})
		return
	}

	transaction, err := h.transactionService.GetTransactionByID(c.Request.Context(), transactionID, accessToken)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "transaction not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "access denied" {
			statusCode = http.StatusForbidden
		}

		c.JSON(statusCode, ErrorResponse{
			Error:   "TRANSACTION_NOT_FOUND",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// GetUserTransactions возвращает список транзакций пользователя
// @Summary Get user transactions
// @Description Get paginated list of transactions for the authenticated user
// @Tags transactions
// @Produce json
// @Param user_id path string true "User ID"
// @Param limit query int false "Limit (default: 10)"
// @Param offset query int false "Offset (default: 0)"
// @Security BearerAuth
// @Success 200 {array} model.Transaction
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /users/{user_id}/transactions [get]
func (h *TransactionHandler) GetUserTransactions(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "VALIDATION_ERROR",
			Message: "User ID is required",
		})
		return
	}

	limit := 10
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	accessToken, err := security.ExtractTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "AUTH_ERROR",
			Message: err.Error(),
		})
		return
	}

	transactions, err := h.transactionService.GetUserTransactions(c.Request.Context(), userID, limit, offset, accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "SERVER_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// GetTransactionStats возвращает статистику транзакций пользователя
// @Summary Get transaction statistics
// @Description Get transaction statistics for the authenticated user
// @Tags transactions
// @Produce json
// @Param user_id path string true "User ID"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /users/{user_id}/transactions/stats [get]
func (h *TransactionHandler) GetTransactionStats(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "VALIDATION_ERROR",
			Message: "User ID is required",
		})
		return
	}

	accessToken, err := security.ExtractTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "AUTH_ERROR",
			Message: err.Error(),
		})
		return
	}

	stats, err := h.transactionService.GetTransactionStats(c.Request.Context(), userID, accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "SERVER_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// HealthCheck проверка здоровья сервиса
// @Summary Health check
// @Description Check if the service is healthy
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (h *TransactionHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "transaction-service",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// BatchCreateTransactions создает несколько транзакций
// @Summary Create multiple transactions
// @Description Process multiple financial transactions in batch
// @Tags transactions
// @Accept json
// @Produce json
// @Param request body []CreateTransactionRequest true "Batch transaction requests"
// @Security BearerAuth
// @Success 201 {object} service.BatchTransactionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /transactions/batch [post]
func (h *TransactionHandler) BatchCreateTransactions(c *gin.Context) {
	var requests []CreateTransactionRequest
	if err := c.ShouldBindJSON(&requests); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "VALIDATION_ERROR",
			Message: err.Error(),
		})
		return
	}

	accessToken, err := security.ExtractTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "AUTH_ERROR",
			Message: err.Error(),
		})
		return
	}

	claims, exists := c.Get("userClaims")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "AUTH_ERROR",
			Message: "User claims not found",
		})
		return
	}

	jwtClaims, ok := claims.(*security.JWTClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "AUTH_ERROR",
			Message: "Invalid token claims",
		})
		return
	}

	serviceRequests := make([]*service.ProcessTransactionRequest, len(requests))
	for i, req := range requests {
		serviceRequests[i] = &service.ProcessTransactionRequest{
			UserID:      jwtClaims.UserID,
			AccountID:   req.AccountID,
			Amount:      req.Amount,
			Currency:    req.Currency,
			Type:        req.Type,
			Description: req.Description,
			Reference:   req.Reference,
			AccessToken: accessToken,
		}
	}

	result, err := h.transactionService.ProcessBatchTransactions(c.Request.Context(), serviceRequests)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "BATCH_TRANSACTION_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// ErrorResponse стандартный формат ошибки
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
