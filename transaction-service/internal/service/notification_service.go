package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"transaction-service/internal/model"
	"transaction-service/internal/repository/postgres"
	"transaction-service/internal/repository/redis"
)

type NotificationService struct {
	notificationRepo *postgres.NotificationRepository
	cacheRepo        *redis.CacheRepository
	emailService     *EmailService
	smsService       *SMSService
	pushService      *PushNotificationService
}

type TransactionNotificationRequest struct {
	UserID        string                 `json:"user_id"`
	TransactionID string                 `json:"transaction_id"`
	Type          model.NotificationType `json:"type"`
	Title         string                 `json:"title"`
	Message       string                 `json:"message"`
	Data          map[string]interface{} `json:"data"`
	Channels      []string               `json:"channels"` // EMAIL, SMS, PUSH
}

func NewNotificationService(
	notificationRepo *postgres.NotificationRepository,
	cacheRepo *redis.CacheRepository,
	emailService *EmailService,
	smsService *SMSService,
	pushService *PushNotificationService,
) *NotificationService {
	return &NotificationService{
		notificationRepo: notificationRepo,
		cacheRepo:        cacheRepo,
		emailService:     emailService,
		smsService:       smsService,
		pushService:      pushService,
	}
}

func (s *NotificationService) CreateAndSendNotification(ctx context.Context, req *model.NotificationRequest) (*model.Notification, error) {
	// 1. Создаем уведомление в БД
	notification := &model.Notification{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		Type:      req.Type,
		Title:     req.Title,
		Message:   req.Message,
		Data:      req.Data,
		Status:    model.NotificationStatusPending,
		Channel:   req.Channel,
		CreatedAt: time.Now(),
	}

	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	// 2. Отправляем уведомление через выбранный канал
	var sendErr error
	switch req.Channel {
	case "EMAIL":
		sendErr = s.emailService.Send(ctx, notification)
	case "SMS":
		sendErr = s.smsService.Send(ctx, notification)
	case "PUSH":
		sendErr = s.pushService.Send(ctx, notification)
	default:
		sendErr = fmt.Errorf("unsupported notification channel: %s", req.Channel)
	}

	// 3. Обновляем статус уведомления
	if sendErr != nil {
		notification.Status = model.NotificationStatusFailed
		if updateErr := s.notificationRepo.UpdateStatus(ctx, notification.ID, model.NotificationStatusFailed); updateErr != nil {
			return nil, fmt.Errorf("failed to send notification: %w, and failed to update status: %v", sendErr, updateErr)
		}
		return nil, fmt.Errorf("failed to send notification: %w", sendErr)
	}

	// Успешная отправка
	now := time.Now()
	notification.Status = model.NotificationStatusSent
	notification.SentAt = &now

	if err := s.notificationRepo.MarkAsSent(ctx, notification.ID); err != nil {
		return nil, fmt.Errorf("failed to mark notification as sent: %w", err)
	}

	// 4. Инвалидируем кэш
	_ = s.cacheRepo.InvalidateUserNotifications(ctx, req.UserID)

	return notification, nil
}

func (s *NotificationService) SendTransactionNotification(ctx context.Context, req *TransactionNotificationRequest) error {
	// Отправляем уведомления через все выбранные каналы
	for _, channel := range req.Channels {
		notificationReq := &model.NotificationRequest{
			UserID:  req.UserID,
			Type:    req.Type,
			Title:   req.Title,
			Message: req.Message,
			Data:    req.Data,
			Channel: channel,
		}

		_, err := s.CreateAndSendNotification(ctx, notificationReq)
		if err != nil {
			// Логируем ошибку, но продолжаем отправку через другие каналы
			fmt.Printf("Failed to send %s notification: %v\n", channel, err)
		}
	}

	return nil
}

func (s *NotificationService) GetUserNotifications(ctx context.Context, userID string, limit, offset int) ([]*model.Notification, error) {
	// Сначала проверяем кэш
	cacheKey := fmt.Sprintf("user:%s:notifications:%d:%d", userID, limit, offset)
	var cachedNotifications []*model.Notification
	found, err := s.cacheRepo.Get(ctx, cacheKey, &cachedNotifications)
	if err != nil {
		return nil, fmt.Errorf("cache error: %w", err)
	}
	if found {
		return cachedNotifications, nil
	}

	// Если нет в кэше, ищем в БД
	notifications, err := s.notificationRepo.FindByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кэш
	_ = s.cacheRepo.Set(ctx, cacheKey, notifications)

	return notifications, nil
}

func (s *NotificationService) MarkNotificationAsSent(ctx context.Context, notificationID string) error {
	return s.notificationRepo.MarkAsSent(ctx, notificationID)
}

func (s *NotificationService) GetNotificationStats(ctx context.Context, userID string) (map[string]interface{}, error) {
	// Сначала проверяем кэш
	cacheKey := fmt.Sprintf("user:%s:notification_stats", userID)
	var cachedStats map[string]interface{}
	found, err := s.cacheRepo.Get(ctx, cacheKey, &cachedStats)
	if err != nil {
		return nil, fmt.Errorf("cache error: %w", err)
	}
	if found {
		return cachedStats, nil
	}

	// Если нет в кэше, вычисляем статистику
	stats, err := s.notificationRepo.GetUserNotificationStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кэш с меньшим TTL
	_ = s.cacheRepo.SetWithCustomTTL(ctx, cacheKey, stats, 10*time.Minute)

	return stats, nil
}

func (s *NotificationService) ProcessPendingNotifications(ctx context.Context) (*BatchNotificationResponse, error) {
	// Получаем все pending уведомления
	pendingNotifications, err := s.notificationRepo.FindPendingNotifications(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending notifications: %w", err)
	}

	response := &BatchNotificationResponse{
		Successful: make([]*NotificationResult, 0),
		Failed:     make([]*NotificationResult, 0),
	}

	// Обрабатываем каждое уведомление
	for _, notification := range pendingNotifications {
		req := &model.NotificationRequest{
			UserID:  notification.UserID,
			Type:    notification.Type,
			Title:   notification.Title,
			Message: notification.Message,
			Data:    notification.Data,
			Channel: notification.Channel,
		}

		result, err := s.CreateAndSendNotification(ctx, req)
		notificationResult := &NotificationResult{
			NotificationID: notification.ID,
			UserID:         notification.UserID,
			Channel:        notification.Channel,
			Type:           string(notification.Type),
		}

		if err != nil {
			notificationResult.Error = err.Error()
			response.Failed = append(response.Failed, notificationResult)
		} else {
			notificationResult.Success = true
			response.Successful = append(response.Successful, notificationResult)
		}
	}

	response.Total = len(pendingNotifications)
	response.SuccessCount = len(response.Successful)
	response.FailureCount = len(response.Failed)

	return response, nil
}

// Вспомогательные сервисы для разных каналов уведомлений

type EmailService struct {
	// Конфигурация email сервиса
}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (s *EmailService) Send(ctx context.Context, notification *model.Notification) error {
	// Реализация отправки email
	// Интеграция с SMTP сервером или email API
	fmt.Printf("Sending email to user %s: %s - %s\n",
		notification.UserID, notification.Title, notification.Message)
	return nil
}

type SMSService struct {
	// Конфигурация SMS сервиса
}

func NewSMSService() *SMSService {
	return &SMSService{}
}

func (s *SMSService) Send(ctx context.Context, notification *model.Notification) error {
	// Реализация отправки SMS
	// Интеграция с SMS gateway
	fmt.Printf("Sending SMS to user %s: %s\n",
		notification.UserID, notification.Message)
	return nil
}

type PushNotificationService struct {
	// Конфигурация push notification сервиса
}

func NewPushNotificationService() *PushNotificationService {
	return &PushNotificationService{}
}

func (s *PushNotificationService) Send(ctx context.Context, notification *model.Notification) error {
	// Реализация отправки push уведомлений
	// Интеграция с FCM, APNS и т.д.
	fmt.Printf("Sending push notification to user %s: %s - %s\n",
		notification.UserID, notification.Title, notification.Message)
	return nil
}

// Вспомогательные структуры

type BatchNotificationResponse struct {
	Successful   []*NotificationResult `json:"successful"`
	Failed       []*NotificationResult `json:"failed"`
	Total        int                   `json:"total"`
	SuccessCount int                   `json:"success_count"`
	FailureCount int                   `json:"failure_count"`
}

type NotificationResult struct {
	NotificationID string `json:"notification_id"`
	UserID         string `json:"user_id"`
	Channel        string `json:"channel"`
	Type           string `json:"type"`
	Success        bool   `json:"success"`
	Error          string `json:"error,omitempty"`
}
