package handler

import (
	"net/http"
	"strconv"

	"transaction-service/internal/model"
	"transaction-service/internal/service"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notificationService *service.NotificationService
}

func NewNotificationHandler(notificationService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

type CreateNotificationRequest struct {
	UserID  string                 `json:"user_id" binding:"required"`
	Type    model.NotificationType `json:"type" binding:"required"`
	Title   string                 `json:"title" binding:"required"`
	Message string                 `json:"message" binding:"required"`
	Data    map[string]interface{} `json:"data"`
	Channel string                 `json:"channel" binding:"required"` // EMAIL, SMS, PUSH
}

// CreateNotification создает новое уведомление
// @Summary Create notification
// @Description Create and send a notification to user
// @Tags notifications
// @Accept json
// @Produce json
// @Param request body CreateNotificationRequest true "Notification request"
// @Security BearerAuth
// @Success 201 {object} model.Notification
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /notifications [post]
func (h *NotificationHandler) CreateNotification(c *gin.Context) {
	var req CreateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "VALIDATION_ERROR",
			Message: err.Error(),
		})
		return
	}

	notificationReq := &model.NotificationRequest{
		UserID:  req.UserID,
		Type:    req.Type,
		Title:   req.Title,
		Message: req.Message,
		Data:    req.Data,
		Channel: req.Channel,
	}

	notification, err := h.notificationService.CreateAndSendNotification(c.Request.Context(), notificationReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "NOTIFICATION_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, notification)
}

// GetUserNotifications возвращает уведомления пользователя
// @Summary Get user notifications
// @Description Get paginated list of notifications for user
// @Tags notifications
// @Produce json
// @Param user_id path string true "User ID"
// @Param limit query int false "Limit (default: 10)"
// @Param offset query int false "Offset (default: 0)"
// @Security BearerAuth
// @Success 200 {array} model.Notification
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /users/{user_id}/notifications [get]
func (h *NotificationHandler) GetUserNotifications(c *gin.Context) {
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

	notifications, err := h.notificationService.GetUserNotifications(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "SERVER_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// MarkAsSent помечает уведомление как отправленное
// @Summary Mark notification as sent
// @Description Mark notification as sent by ID
// @Tags notifications
// @Produce json
// @Param id path string true "Notification ID"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /notifications/{id}/mark-sent [put]
func (h *NotificationHandler) MarkAsSent(c *gin.Context) {
	notificationID := c.Param("id")
	if notificationID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "VALIDATION_ERROR",
			Message: "Notification ID is required",
		})
		return
	}

	err := h.notificationService.MarkNotificationAsSent(c.Request.Context(), notificationID)
	if err != nil {
		if err.Error() == "notification not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "NOT_FOUND",
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "SERVER_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Notification marked as sent",
		"id":      notificationID,
	})
}

// SendTransactionNotification отправляет уведомление о транзакции
// @Summary Send transaction notification
// @Description Send notification for a transaction event
// @Tags notifications
// @Accept json
// @Produce json
// @Param request body service.TransactionNotificationRequest true "Transaction notification request"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Router /notifications/transaction [post]
func (h *NotificationHandler) SendTransactionNotification(c *gin.Context) {
	var req service.TransactionNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "VALIDATION_ERROR",
			Message: err.Error(),
		})
		return
	}

	err := h.notificationService.SendTransactionNotification(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "NOTIFICATION_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction notification sent successfully",
	})
}

// GetNotificationStats возвращает статистику уведомлений
// @Summary Get notification statistics
// @Description Get notification statistics for user
// @Tags notifications
// @Produce json
// @Param user_id path string true "User ID"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Router /users/{user_id}/notifications/stats [get]
func (h *NotificationHandler) GetNotificationStats(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "VALIDATION_ERROR",
			Message: "User ID is required",
		})
		return
	}

	stats, err := h.notificationService.GetNotificationStats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "SERVER_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}
