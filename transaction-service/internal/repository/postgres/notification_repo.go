package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"transaction-service/internal/model"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(ctx context.Context, notification *model.Notification) error {
	query := `
		INSERT INTO notifications (
			id, user_id, type, title, message, data, status, channel, created_at, sent_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	dataJSON, err := json.Marshal(notification.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal notification data: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query,
		notification.ID,
		notification.UserID,
		string(notification.Type),
		notification.Title,
		notification.Message,
		dataJSON,
		string(notification.Status),
		notification.Channel,
		notification.CreatedAt,
		notification.SentAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	return nil
}

func (r *NotificationRepository) FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, data, status, channel, created_at, sent_at
		FROM notifications 
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query notifications: %w", err)
	}
	defer rows.Close()

	var notifications []*model.Notification
	for rows.Next() {
		var notification model.Notification
		var dataJSON string
		var typeStr, statusStr string
		var sentAt sql.NullTime

		err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&typeStr,
			&notification.Title,
			&notification.Message,
			&dataJSON,
			&statusStr,
			&notification.Channel,
			&notification.CreatedAt,
			&sentAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}

		notification.Type = model.NotificationType(typeStr)
		notification.Status = model.NotificationStatus(statusStr)

		if sentAt.Valid {
			notification.SentAt = &sentAt.Time
		}

		if err := json.Unmarshal([]byte(dataJSON), &notification.Data); err != nil {
			return nil, fmt.Errorf("failed to unmarshal notification data: %w", err)
		}

		notifications = append(notifications, &notification)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return notifications, nil
}

func (r *NotificationRepository) UpdateStatus(ctx context.Context, id string, status model.NotificationStatus) error {
	query := `
		UPDATE notifications 
		SET status = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, string(status), id)
	if err != nil {
		return fmt.Errorf("failed to update notification status: %w", err)
	}

	return nil
}

func (r *NotificationRepository) MarkAsSent(ctx context.Context, id string) error {
	query := `
		UPDATE notifications 
		SET status = $1, sent_at = $2
		WHERE id = $3
	`

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, string(model.NotificationStatusSent), now, id)
	if err != nil {
		return fmt.Errorf("failed to mark notification as sent: %w", err)
	}

	return nil
}

func (r *NotificationRepository) FindPendingNotifications(ctx context.Context) ([]*model.Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, data, status, channel, created_at, sent_at
		FROM notifications 
		WHERE status = 'PENDING'
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query pending notifications: %w", err)
	}
	defer rows.Close()

	var notifications []*model.Notification
	for rows.Next() {
		var notification model.Notification
		var dataJSON string
		var typeStr, statusStr string
		var sentAt sql.NullTime

		err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&typeStr,
			&notification.Title,
			&notification.Message,
			&dataJSON,
			&statusStr,
			&notification.Channel,
			&notification.CreatedAt,
			&sentAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}

		notification.Type = model.NotificationType(typeStr)
		notification.Status = model.NotificationStatus(statusStr)

		if sentAt.Valid {
			notification.SentAt = &sentAt.Time
		}

		if err := json.Unmarshal([]byte(dataJSON), &notification.Data); err != nil {
			return nil, fmt.Errorf("failed to unmarshal notification data: %w", err)
		}

		notifications = append(notifications, &notification)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return notifications, nil
}

func (r *NotificationRepository) GetUserNotificationStats(ctx context.Context, userID string) (map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(*) as total_count,
			COUNT(CASE WHEN status = 'SENT' THEN 1 END) as sent_count,
			COUNT(CASE WHEN status = 'PENDING' THEN 1 END) as pending_count,
			COUNT(CASE WHEN status = 'FAILED' THEN 1 END) as failed_count,
			COUNT(CASE WHEN channel = 'EMAIL' THEN 1 END) as email_count,
			COUNT(CASE WHEN channel = 'SMS' THEN 1 END) as sms_count,
			COUNT(CASE WHEN channel = 'PUSH' THEN 1 END) as push_count
		FROM notifications 
		WHERE user_id = $1
	`

	stats := make(map[string]interface{})
	var totalCount, sentCount, pendingCount, failedCount, emailCount, smsCount, pushCount int

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&totalCount,
		&sentCount,
		&pendingCount,
		&failedCount,
		&emailCount,
		&smsCount,
		&pushCount,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get notification stats: %w", err)
	}

	stats["total_count"] = totalCount
	stats["sent_count"] = sentCount
	stats["pending_count"] = pendingCount
	stats["failed_count"] = failedCount
	stats["email_count"] = emailCount
	stats["sms_count"] = smsCount
	stats["push_count"] = pushCount
	stats["success_rate"] = 0.0

	if totalCount > 0 {
		stats["success_rate"] = float64(sentCount) / float64(totalCount) * 100
	}

	return stats, nil
}

func (r *NotificationRepository) FindByID(ctx context.Context, id string) (*model.Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, data, status, channel, created_at, sent_at
		FROM notifications 
		WHERE id = $1
	`

	var notification model.Notification
	var dataJSON string
	var typeStr, statusStr string
	var sentAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&notification.ID,
		&notification.UserID,
		&typeStr,
		&notification.Title,
		&notification.Message,
		&dataJSON,
		&statusStr,
		&notification.Channel,
		&notification.CreatedAt,
		&sentAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("notification not found: %s", id)
		}
		return nil, fmt.Errorf("failed to find notification: %w", err)
	}

	notification.Type = model.NotificationType(typeStr)
	notification.Status = model.NotificationStatus(statusStr)

	if sentAt.Valid {
		notification.SentAt = &sentAt.Time
	}

	if err := json.Unmarshal([]byte(dataJSON), &notification.Data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal notification data: %w", err)
	}

	return &notification, nil
}
