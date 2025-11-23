package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"transaction-service/internal/client"
	"transaction-service/internal/model"
	"transaction-service/internal/repository/postgres"
	"transaction-service/internal/repository/redis"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInvalidAmount       = errors.New("invalid amount")
	ErrAccountNotFound     = errors.New("account not found")
	ErrTransactionFailed   = errors.New("transaction failed")
)

type TransactionService struct {
	transactionRepo *postgres.TransactionRepository
	cacheRepo       *redis.CacheRepository
	authClient      *client.AuthClient
	depositClient   *client.DepositClient
}

func NewTransactionService(
	transactionRepo *postgres.TransactionRepository,
	cacheRepo *redis.CacheRepository,
	authClient *client.AuthClient,
	depositClient *client.DepositClient,
) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		cacheRepo:       cacheRepo,
		authClient:      authClient,
		depositClient:   depositClient,
	}
}

type ProcessTransactionRequest struct {
	UserID      string                `json:"user_id"`
	AccountID   string                `json:"account_id"`
	Amount      float64               `json:"amount"`
	Currency    string                `json:"currency"`
	Type        model.TransactionType `json:"type"`
	Description string                `json:"description"`
	Reference   string                `json:"reference"`
	AccessToken string                `json:"-"`
}

type ProcessTransactionResponse struct {
	Transaction *model.Transaction `json:"transaction"`
	Message     string             `json:"message"`
	NewBalance  float64            `json:"new_balance,omitempty"`
}

func (s *TransactionService) ProcessTransaction(ctx context.Context, req *ProcessTransactionRequest) (*ProcessTransactionResponse, error) {
	// 1. Валидация входных данных
	if err := s.validateTransactionRequest(req); err != nil {
		return nil, err
	}

	// 2. Проверка аутентификации пользователя
	userInfo, err := s.authClient.ValidateToken(ctx, req.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	if userInfo.UserID != req.UserID {
		return nil, errors.New("user ID mismatch")
	}

	// 3. Проверка баланса для операций списания
	if req.Type == model.TransactionTypeWithdrawal {
		balance, err := s.depositClient.GetBalance(ctx, req.AccountID, req.AccessToken)
		if err != nil {
			return nil, fmt.Errorf("failed to get balance: %w", err)
		}

		if balance.Available < req.Amount {
			return nil, ErrInsufficientBalance
		}
	}

	// 4. Создание транзакции
	transactionReq := &model.TransactionRequest{
		UserID:      req.UserID,
		AccountID:   req.AccountID,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Type:        req.Type,
		Description: req.Description,
		Reference:   req.Reference,
	}

	transaction := model.NewTransaction(transactionReq)

	// 5. Сохранение транзакции в PostgreSQL
	if err := s.transactionRepo.Create(ctx, transaction); err != nil {
		return nil, fmt.Errorf("failed to save transaction: %w", err)
	}

	// 6. Выполнение операции в deposit-service
	depositReq := &client.TransactionRequest{
		FromAccountID: req.AccountID,
		Amount:        req.Amount,
		Currency:      req.Currency,
		Type:          string(req.Type),
		Description:   req.Description,
		Reference:     req.Reference,
	}

	depositResp, err := s.depositClient.ProcessTransaction(ctx, depositReq, req.AccessToken)
	if err != nil {
		// Обновляем статус транзакции на FAILED
		_ = s.transactionRepo.UpdateStatus(ctx, transaction.ID, model.TransactionStatusFailed)
		return nil, fmt.Errorf("failed to process transaction in deposit service: %w", err)
	}

	// 7. Обновление статуса транзакции на COMPLETED
	if err := s.transactionRepo.UpdateStatus(ctx, transaction.ID, model.TransactionStatusCompleted); err != nil {
		return nil, fmt.Errorf("failed to update transaction status: %w", err)
	}

	transaction.Status = model.TransactionStatusCompleted

	// 8. Инвалидация кэша
	_ = s.cacheRepo.InvalidateUserTransactions(ctx, req.UserID)
	_ = s.cacheRepo.InvalidateTransaction(ctx, transaction.ID)

	response := &ProcessTransactionResponse{
		Transaction: transaction,
		Message:     "Transaction processed successfully",
		NewBalance:  depositResp.NewBalance,
	}

	return response, nil
}

func (s *TransactionService) GetTransactionByID(ctx context.Context, transactionID, accessToken string) (*model.Transaction, error) {
	// Сначала проверяем кэш
	var cachedTransaction model.Transaction
	found, err := s.cacheRepo.GetCachedTransaction(ctx, transactionID, &cachedTransaction)
	if err != nil {
		return nil, fmt.Errorf("cache error: %w", err)
	}
	if found {
		return &cachedTransaction, nil
	}

	// Если нет в кэше, ищем в БД
	transaction, err := s.transactionRepo.FindByID(ctx, transactionID)
	if err != nil {
		return nil, err
	}

	// Проверяем права доступа
	userInfo, err := s.authClient.ValidateToken(ctx, accessToken)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	if userInfo.UserID != transaction.UserID && !s.hasAdminRole(userInfo.Roles) {
		return nil, errors.New("access denied")
	}

	// Сохраняем в кэш
	_ = s.cacheRepo.CacheTransaction(ctx, transactionID, transaction)

	return transaction, nil
}

func (s *TransactionService) GetUserTransactions(ctx context.Context, userID string, limit, offset int, accessToken string) ([]*model.Transaction, error) {
	// Проверяем права доступа
	userInfo, err := s.authClient.ValidateToken(ctx, accessToken)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	if userInfo.UserID != userID && !s.hasAdminRole(userInfo.Roles) {
		return nil, errors.New("access denied")
	}

	// Сначала проверяем кэш
	cacheKey := fmt.Sprintf("user:%s:transactions:%d:%d", userID, limit, offset)
	var cachedTransactions []*model.Transaction
	found, err := s.cacheRepo.Get(ctx, cacheKey, &cachedTransactions)
	if err != nil {
		return nil, fmt.Errorf("cache error: %w", err)
	}
	if found {
		return cachedTransactions, nil
	}

	// Если нет в кэше, ищем в БД
	transactions, err := s.transactionRepo.FindByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кэш
	_ = s.cacheRepo.Set(ctx, cacheKey, transactions)

	return transactions, nil
}

func (s *TransactionService) GetTransactionStats(ctx context.Context, userID, accessToken string) (map[string]interface{}, error) {
	// Проверяем права доступа
	userInfo, err := s.authClient.ValidateToken(ctx, accessToken)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	if userInfo.UserID != userID && !s.hasAdminRole(userInfo.Roles) {
		return nil, errors.New("access denied")
	}

	// Сначала проверяем кэш
	cacheKey := fmt.Sprintf("user:%s:stats", userID)
	var cachedStats map[string]interface{}
	found, err := s.cacheRepo.Get(ctx, cacheKey, &cachedStats)
	if err != nil {
		return nil, fmt.Errorf("cache error: %w", err)
	}
	if found {
		return cachedStats, nil
	}

	// Если нет в кэше, вычисляем статистику
	stats, err := s.transactionRepo.GetUserTransactionStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кэш с меньшим TTL
	_ = s.cacheRepo.SetWithCustomTTL(ctx, cacheKey, stats, 5*time.Minute)

	return stats, nil
}

func (s *TransactionService) ProcessBatchTransactions(ctx context.Context, requests []*ProcessTransactionRequest) (*BatchTransactionResponse, error) {
	results := &BatchTransactionResponse{
		Successful: make([]*ProcessTransactionResponse, 0),
		Failed:     make([]*FailedTransaction, 0),
	}

	for _, req := range requests {
		result, err := s.ProcessTransaction(ctx, req)
		if err != nil {
			results.Failed = append(results.Failed, &FailedTransaction{
				Request: req,
				Error:   err.Error(),
			})
		} else {
			results.Successful = append(results.Successful, result)
		}
	}

	results.Total = len(requests)
	results.SuccessCount = len(results.Successful)
	results.FailureCount = len(results.Failed)

	return results, nil
}

func (s *TransactionService) validateTransactionRequest(req *ProcessTransactionRequest) error {
	if req.Amount <= 0 {
		return ErrInvalidAmount
	}

	if req.UserID == "" {
		return errors.New("user ID is required")
	}

	if req.AccountID == "" {
		return errors.New("account ID is required")
	}

	if req.Currency == "" {
		return errors.New("currency is required")
	}

	if req.AccessToken == "" {
		return errors.New("access token is required")
	}

	return nil
}

func (s *TransactionService) hasAdminRole(roles []string) bool {
	for _, role := range roles {
		if role == "ROLE_ADMIN" || role == "ADMIN" {
			return true
		}
	}
	return false
}

// Вспомогательные структуры для batch processing
type BatchTransactionResponse struct {
	Successful   []*ProcessTransactionResponse `json:"successful"`
	Failed       []*FailedTransaction          `json:"failed"`
	Total        int                           `json:"total"`
	SuccessCount int                           `json:"success_count"`
	FailureCount int                           `json:"failure_count"`
}

type FailedTransaction struct {
	Request *ProcessTransactionRequest `json:"request"`
	Error   string                     `json:"error"`
}
