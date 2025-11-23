package model

import "time"

type NotificationType string

const (
	NotificationTypeTransaction NotificationType = "TRANSACTION"
	NotificationTypeInterest    NotificationType = "INTEREST"
	NotificationTypeSystem      NotificationType = "SYSTEM"
)

type NotificationStatus string

const (
	NotificationStatusPending NotificationStatus = "PENDING"
	NotificationStatusSent    NotificationStatus = "SENT"
	NotificationStatusFailed  NotificationStatus = "FAILED"
)

type Notification struct {
	ID        string                 `json:"id" db:"id"`
	UserID    string                 `json:"user_id" db:"user_id"`
	Type      NotificationType       `json:"type" db:"type"`
	Title     string                 `json:"title" db:"title"`
	Message   string                 `json:"message" db:"message"`
	Data      map[string]interface{} `json:"data" db:"data"`
	Status    NotificationStatus     `json:"status" db:"status"`
	Channel   string                 `json:"channel" db:"channel"` // EMAIL, SMS, PUSH
	CreatedAt time.Time              `json:"created_at" db:"created_at"`
	SentAt    *time.Time             `json:"sent_at" db:"sent_at"`
}

type NotificationRequest struct {
	UserID  string                 `json:"user_id" binding:"required"`
	Type    NotificationType       `json:"type" binding:"required"`
	Title   string                 `json:"title" binding:"required"`
	Message string                 `json:"message" binding:"required"`
	Data    map[string]interface{} `json:"data"`
	Channel string                 `json:"channel" binding:"required"`
}
