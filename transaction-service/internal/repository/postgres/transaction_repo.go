package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"transaction-service/internal/model"

	_ "github.com/lib/pq"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, transaction *model.Transaction) error {
	query := `
		INSERT INTO transactions (
			id, user_id, account_id, amount, currency, type, status, 
			description, reference, metadata, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	metadataJSON, err := json.Marshal(transaction.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query,
		transaction.ID,
		transaction.UserID,
		transaction.AccountID,
		transaction.Amount,
		transaction.Currency,
		string(transaction.Type),
		string(transaction.Status),
		transaction.Description,
		transaction.Reference,
		metadataJSON,
		transaction.CreatedAt,
		transaction.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}

func (r *TransactionRepository) FindByID(ctx context.Context, id string) (*model.Transaction, error) {
	query := `
		SELECT id, user_id, account_id, amount, currency, type, status,
		       description, reference, metadata, created_at, updated_at
		FROM transactions 
		WHERE id = $1
	`

	var transaction model.Transaction
	var metadataJSON string
	var typeStr, statusStr string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.AccountID,
		&transaction.Amount,
		&transaction.Currency,
		&typeStr,
		&statusStr,
		&transaction.Description,
		&transaction.Reference,
		&metadataJSON,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction not found: %s", id)
		}
		return nil, fmt.Errorf("failed to find transaction: %w", err)
	}

	transaction.Type = model.TransactionType(typeStr)
	transaction.Status = model.TransactionStatus(statusStr)

	if err := json.Unmarshal([]byte(metadataJSON), &transaction.Metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return &transaction, nil
}

func (r *TransactionRepository) FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.Transaction, error) {
	query := `
		SELECT id, user_id, account_id, amount, currency, type, status,
		       description, reference, metadata, created_at, updated_at
		FROM transactions 
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*model.Transaction
	for rows.Next() {
		var transaction model.Transaction
		var metadataJSON string
		var typeStr, statusStr string

		err := rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.AccountID,
			&transaction.Amount,
			&transaction.Currency,
			&typeStr,
			&statusStr,
			&transaction.Description,
			&transaction.Reference,
			&metadataJSON,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}

		transaction.Type = model.TransactionType(typeStr)
		transaction.Status = model.TransactionStatus(statusStr)

		if err := json.Unmarshal([]byte(metadataJSON), &transaction.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		transactions = append(transactions, &transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return transactions, nil
}

func (r *TransactionRepository) UpdateStatus(ctx context.Context, id string, status model.TransactionStatus) error {
	query := `
		UPDATE transactions 
		SET status = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.ExecContext(ctx, query, string(status), time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}

	return nil
}

func (r *TransactionRepository) FindByStatus(ctx context.Context, status model.TransactionStatus) ([]*model.Transaction, error) {
	query := `
		SELECT id, user_id, account_id, amount, currency, type, status,
		       description, reference, metadata, created_at, updated_at
		FROM transactions 
		WHERE status = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, string(status))
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions by status: %w", err)
	}
	defer rows.Close()

	var transactions []*model.Transaction
	for rows.Next() {
		var transaction model.Transaction
		var metadataJSON string
		var typeStr, statusStr string

		err := rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.AccountID,
			&transaction.Amount,
			&transaction.Currency,
			&typeStr,
			&statusStr,
			&transaction.Description,
			&transaction.Reference,
			&metadataJSON,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}

		transaction.Type = model.TransactionType(typeStr)
		transaction.Status = model.TransactionStatus(statusStr)

		if err := json.Unmarshal([]byte(metadataJSON), &transaction.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		transactions = append(transactions, &transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return transactions, nil
}

func (r *TransactionRepository) GetUserTransactionStats(ctx context.Context, userID string) (map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(*) as total_count,
			COUNT(CASE WHEN status = 'COMPLETED' THEN 1 END) as completed_count,
			COUNT(CASE WHEN status = 'FAILED' THEN 1 END) as failed_count,
			COALESCE(SUM(CASE WHEN type = 'DEPOSIT' AND status = 'COMPLETED' THEN amount ELSE 0 END), 0) as total_deposits,
			COALESCE(SUM(CASE WHEN type = 'WITHDRAWAL' AND status = 'COMPLETED' THEN amount ELSE 0 END), 0) as total_withdrawals
		FROM transactions 
		WHERE user_id = $1
	`

	stats := make(map[string]interface{})
	var totalCount, completedCount, failedCount int
	var totalDeposits, totalWithdrawals float64

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&totalCount,
		&completedCount,
		&failedCount,
		&totalDeposits,
		&totalWithdrawals,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get user transaction stats: %w", err)
	}

	stats["total_count"] = totalCount
	stats["completed_count"] = completedCount
	stats["failed_count"] = failedCount
	stats["total_deposits"] = totalDeposits
	stats["total_withdrawals"] = totalWithdrawals
	stats["net_flow"] = totalDeposits - totalWithdrawals

	return stats, nil
}
