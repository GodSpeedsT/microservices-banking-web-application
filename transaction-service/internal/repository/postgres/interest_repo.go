package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type InterestRate struct {
	ID            string    `json:"id"`
	AccountType   string    `json:"account_type"`
	Rate          float64   `json:"rate"` // годовая процентная ставка
	EffectiveFrom time.Time `json:"effective_from"`
	EffectiveTo   time.Time `json:"effective_to"`
	CreatedAt     time.Time `json:"created_at"`
}

type InterestAccrual struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	AccountID string    `json:"account_id"`
	Period    string    `json:"period"` // YYYY-MM
	Principal float64   `json:"principal"`
	Interest  float64   `json:"interest"`
	Rate      float64   `json:"rate"`
	Status    string    `json:"status"` // PENDING, APPLIED, FAILED
	AppliedAt time.Time `json:"applied_at"`
	CreatedAt time.Time `json:"created_at"`
}

type InterestRepository struct {
	db *sql.DB
}

func NewInterestRepository(db *sql.DB) *InterestRepository {
	return &InterestRepository{db: db}
}

func (r *InterestRepository) GetCurrentInterestRate(ctx context.Context, accountType string) (*InterestRate, error) {
	query := `
		SELECT id, account_type, rate, effective_from, effective_to, created_at
		FROM interest_rates 
		WHERE account_type = $1 
		AND effective_from <= $2 
		AND (effective_to IS NULL OR effective_to >= $2)
		ORDER BY effective_from DESC
		LIMIT 1
	`

	var rate InterestRate
	err := r.db.QueryRowContext(ctx, query, accountType, time.Now()).Scan(
		&rate.ID,
		&rate.AccountType,
		&rate.Rate,
		&rate.EffectiveFrom,
		&rate.EffectiveTo,
		&rate.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no interest rate found for account type: %s", accountType)
		}
		return nil, fmt.Errorf("failed to get interest rate: %w", err)
	}

	return &rate, nil
}

func (r *InterestRepository) CreateInterestAccrual(ctx context.Context, accrual *InterestAccrual) error {
	query := `
		INSERT INTO interest_accruals (
			id, user_id, account_id, period, principal, interest, rate, status, applied_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.ExecContext(ctx, query,
		accrual.ID,
		accrual.UserID,
		accrual.AccountID,
		accrual.Period,
		accrual.Principal,
		accrual.Interest,
		accrual.Rate,
		accrual.Status,
		accrual.AppliedAt,
		accrual.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create interest accrual: %w", err)
	}

	return nil
}

func (r *InterestRepository) FindPendingAccruals(ctx context.Context) ([]*InterestAccrual, error) {
	query := `
		SELECT id, user_id, account_id, period, principal, interest, rate, status, applied_at, created_at
		FROM interest_accruals 
		WHERE status = 'PENDING'
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query pending accruals: %w", err)
	}
	defer rows.Close()

	var accruals []*InterestAccrual
	for rows.Next() {
		var accrual InterestAccrual
		err := rows.Scan(
			&accrual.ID,
			&accrual.UserID,
			&accrual.AccountID,
			&accrual.Period,
			&accrual.Principal,
			&accrual.Interest,
			&accrual.Rate,
			&accrual.Status,
			&accrual.AppliedAt,
			&accrual.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan accrual: %w", err)
		}

		accruals = append(accruals, &accrual)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return accruals, nil
}

func (r *InterestRepository) UpdateAccrualStatus(ctx context.Context, id string, status string) error {
	query := `
		UPDATE interest_accruals 
		SET status = $1, applied_at = $2
		WHERE id = $3
	`

	_, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update accrual status: %w", err)
	}

	return nil
}

func (r *InterestRepository) GetAccrualByPeriod(ctx context.Context, userID, accountID, period string) (*InterestAccrual, error) {
	query := `
		SELECT id, user_id, account_id, period, principal, interest, rate, status, applied_at, created_at
		FROM interest_accruals 
		WHERE user_id = $1 AND account_id = $2 AND period = $3
	`

	var accrual InterestAccrual
	err := r.db.QueryRowContext(ctx, query, userID, accountID, period).Scan(
		&accrual.ID,
		&accrual.UserID,
		&accrual.AccountID,
		&accrual.Period,
		&accrual.Principal,
		&accrual.Interest,
		&accrual.Rate,
		&accrual.Status,
		&accrual.AppliedAt,
		&accrual.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get accrual: %w", err)
	}

	return &accrual, nil
}

func (r *InterestRepository) GetUserAccrualHistory(ctx context.Context, userID string, limit, offset int) ([]*InterestAccrual, error) {
	query := `
		SELECT id, user_id, account_id, period, principal, interest, rate, status, applied_at, created_at
		FROM interest_accruals 
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query accrual history: %w", err)
	}
	defer rows.Close()

	var accruals []*InterestAccrual
	for rows.Next() {
		var accrual InterestAccrual
		err := rows.Scan(
			&accrual.ID,
			&accrual.UserID,
			&accrual.AccountID,
			&accrual.Period,
			&accrual.Principal,
			&accrual.Interest,
			&accrual.Rate,
			&accrual.Status,
			&accrual.AppliedAt,
			&accrual.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan accrual: %w", err)
		}

		accruals = append(accruals, &accrual)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return accruals, nil
}
