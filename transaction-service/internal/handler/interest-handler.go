package handler

import (
	"net/http"
	"strconv"

	"transaction-service/internal/service"
	"transaction-service/pkg/security"

	"github.com/gin-gonic/gin"
)

type InterestHandler struct {
	interestService *service.InterestService
}

func NewInterestHandler(interestService *service.InterestService) *InterestHandler {
	return &InterestHandler{
		interestService: interestService,
	}
}

type CalculateInterestRequest struct {
	AccountID string `json:"account_id" binding:"required"`
	Period    string `json:"period" binding:"required"` // YYYY-MM
}

// CalculateInterest рассчитывает проценты для счета
// @Summary Calculate interest
// @Description Calculate monthly interest for an account
// @Tags interest
// @Accept json
// @Produce json
// @Param request body CalculateInterestRequest true "Interest calculation request"
// @Security BearerAuth
// @Success 200 {object} service.CalculateInterestResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /interest/calculate [post]
func (h *InterestHandler) CalculateInterest(c *gin.Context) {
	var req CalculateInterestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "VALIDATION_ERROR",
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

	serviceReq := &service.CalculateInterestRequest{
		UserID:    jwtClaims.UserID,
		AccountID: req.AccountID,
		Period:    req.Period,
	}

	result, err := h.interestService.CalculateMonthlyInterest(c.Request.Context(), serviceReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "INTEREST_CALCULATION_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ApplyInterest применяет рассчитанные проценты к счету
// @Summary Apply interest
// @Description Apply calculated interest to the account
// @Tags interest
// @Accept json
// @Produce json
// @Param request body CalculateInterestRequest true "Interest application request"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /interest/apply [post]
func (h *InterestHandler) ApplyInterest(c *gin.Context) {
	var req CalculateInterestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "VALIDATION_ERROR",
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

	accessToken, err := security.ExtractTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "AUTH_ERROR",
			Message: err.Error(),
		})
		return
	}

	serviceReq := &service.CalculateInterestRequest{
		UserID:    jwtClaims.UserID,
		AccountID: req.AccountID,
		Period:    req.Period,
	}

	err = h.interestService.ApplyInterest(c.Request.Context(), serviceReq, accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "INTEREST_APPLICATION_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Interest applied successfully",
		"user_id":    jwtClaims.UserID,
		"account_id": req.AccountID,
		"period":     req.Period,
	})
}

// ProcessPendingInterest обрабатывает все pending начисления процентов
// @Summary Process pending interest
// @Description Process all pending interest accruals (admin only)
// @Tags interest
// @Produce json
// @Security BearerAuth
// @Success 200 {object} service.BatchInterestResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /interest/process-pending [post]
func (h *InterestHandler) ProcessPendingInterest(c *gin.Context) {
	accessToken, err := security.ExtractTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "AUTH_ERROR",
			Message: err.Error(),
		})
		return
	}

	result, err := h.interestService.ProcessPendingInterestAccruals(c.Request.Context(), accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "INTEREST_PROCESSING_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetInterestHistory возвращает историю начисления процентов
// @Summary Get interest history
// @Description Get interest accrual history for the authenticated user
// @Tags interest
// @Produce json
// @Param user_id path string true "User ID"
// @Param limit query int false "Limit (default: 10)"
// @Param offset query int false "Offset (default: 0)"
// @Security BearerAuth
// @Success 200 {array} postgres.InterestAccrual
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /users/{user_id}/interest/history [get]
func (h *InterestHandler) GetInterestHistory(c *gin.Context) {
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

	history, err := h.interestService.GetInterestAccrualHistory(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "SERVER_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, history)
}
