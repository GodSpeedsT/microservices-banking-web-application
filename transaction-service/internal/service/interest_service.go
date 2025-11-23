package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"transaction-service/internal/client"
	"transaction-service/internal/repository/postgres"
	"transaction-service/internal/repository/redis"
)

type InterestService struct {
	interestRepo    *postgres.InterestRepository
	transactionRepo *postgres.TransactionRepository
	cacheRepo       *redis.CacheRepository
	depositClient   *client.DepositClient
}

func NewInterestService(
	interestRepo *postgres.InterestRepository,
	transactionRepo *postgres.TransactionRepository,
	cacheRepo *redis.CacheRepository,
	depositClient *client.DepositClient,
) *InterestService {
	return &InterestService{
		interestRepo:    interestRepo,
		transactionRepo: transactionRepo,
		cacheRepo:       cacheRepo,
		depositClient:   depositClient,
	}
}

type CalculateInterestRequest struct {
	UserID    string `json:"user_id"`
	AccountID string `json:"account_id"`
	Period    string `json:"period"` // YYYY-MM
}

type CalculateInterestResponse struct {
	UserID    string  `json:"user_id"`
	AccountID string  `json:"account_id"`
	Period    string  `json:"period"`
	Principal float64 `json:"principal"`
	Interest  float64 `json:"interest"`
	Rate      float64 `json:"rate"`
}

func (s *InterestService) CalculateMonthlyInterest(ctx context.Context, req *CalculateInterestRequest) (*CalculateInterestResponse, error) {
	// 1. Получаем текущую процентную ставку
	// В реальном приложении тип счета должен получаться из deposit-service
	accountType := "SAVINGS" // временное значение
	interestRate, err := s.interestRepo.GetCurrentInterestRate(ctx, accountType)
	if err != nil {
		return nil, fmt.Errorf("failed to get interest rate: %w", err)
	}

	// 2. Проверяем, не был ли уже начислен процент за этот период
	existingAccrual, err := s.interestRepo.GetAccrualByPeriod(ctx, req.UserID, req.AccountID, req.Period)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing accrual: %w", err)
	}

	if existingAccrual != nil {
		return nil, fmt.Errorf("interest already calculated for period %s", req.Period)
	}

	// 3. Получаем средний баланс за период (упрощенная логика)
	// В реальном приложении здесь должна быть сложная логика расчета среднего баланса
	averageBalance, err := s.calculateAverageBalance(ctx, req.AccountID, req.Period)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate average balance: %w", err)
	}

	// 4. Рассчитываем проценты (месячная ставка = годовая / 12)
	monthlyRate := interestRate.Rate / 12 / 100
	interest := averageBalance * monthlyRate

	response := &CalculateInterestResponse{
		UserID:    req.UserID,
		AccountID: req.AccountID,
		Period:    req.Period,
		Principal: averageBalance,
		Interest:  interest,
		Rate:      interestRate.Rate,
	}

	return response, nil
}

func (s *InterestService) ApplyInterest(ctx context.Context, req *CalculateInterestRequest, accessToken string) error {
	// 1. Рассчитываем проценты
	calculation, err := s.CalculateMonthlyInterest(ctx, req)
	if err != nil {
		return err
	}

	// 2. Создаем запись о начислении
	accrual := &postgres.InterestAccrual{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		AccountID: req.AccountID,
		Period:    req.Period,
		Principal: calculation.Principal,
		Interest:  calculation.Interest,
		Rate:      calculation.Rate,
		Status:    "PENDING",
		CreatedAt: time.Now(),
	}

	if err := s.interestRepo.CreateInterestAccrual(ctx, accrual); err != nil {
		return fmt.Errorf("failed to create interest accrual: %w", err)
	}

	// 3. Выполняем операцию начисления в deposit-service
	depositReq := &client.TransactionRequest{
		ToAccountID: req.AccountID,
		Amount:      calculation.Interest,
		Currency:    "USD", // должно получаться из счета
		Type:        "INTEREST",
		Description: fmt.Sprintf("Interest for period %s", req.Period),
		Reference:   accrual.ID,
	}

	_, err = s.depositClient.ProcessTransaction(ctx, depositReq, accessToken)
	if err != nil {
		// Обновляем статус на FAILED
		_ = s.interestRepo.UpdateAccrualStatus(ctx, accrual.ID, "FAILED")
		return fmt.Errorf("failed to apply interest in deposit service: %w", err)
	}

	// 4. Обновляем статус на APPLIED
	if err := s.interestRepo.UpdateAccrualStatus(ctx, accrual.ID, "APPLIED"); err != nil {
		return fmt.Errorf("failed to update accrual status: %w", err)
	}

	// 5. Инвалидируем кэш
	_ = s.cacheRepo.InvalidateUserTransactions(ctx, req.UserID)

	return nil
}

func (s *InterestService) ProcessPendingInterestAccruals(ctx context.Context, accessToken string) (*BatchInterestResponse, error) {
	// 1. Получаем все pending начисления
	pendingAccruals, err := s.interestRepo.FindPendingAccruals(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending accruals: %w", err)
	}

	response := &BatchInterestResponse{
		Successful: make([]*InterestAccrualResult, 0),
		Failed:     make([]*InterestAccrualResult, 0),
	}

	// 2. Обрабатываем каждое начисление
	for _, accrual := range pendingAccruals {
		req := &CalculateInterestRequest{
			UserID:    accrual.UserID,
			AccountID: accrual.AccountID,
			Period:    accrual.Period,
		}

		err := s.ApplyInterest(ctx, req, accessToken)
		result := &InterestAccrualResult{
			AccrualID: accrual.ID,
			UserID:    accrual.UserID,
			AccountID: accrual.AccountID,
			Period:    accrual.Period,
			Interest:  accrual.Interest,
		}

		if err != nil {
			result.Error = err.Error()
			response.Failed = append(response.Failed, result)
		} else {
			response.Successful = append(response.Successful, result)
		}
	}

	response.Total = len(pendingAccruals)
	response.SuccessCount = len(response.Successful)
	response.FailureCount = len(response.Failed)

	return response, nil
}

func (s *InterestService) GetInterestAccrualHistory(ctx context.Context, userID string, limit, offset int) ([]*postgres.InterestAccrual, error) {
	return s.interestRepo.GetUserAccrualHistory(ctx, userID, limit, offset)
}

func (s *InterestService) calculateAverageBalance(ctx context.Context, accountID, period string) (float64, error) {

	return 1000.0, nil
}

type BatchInterestResponse struct {
	Successful   []*InterestAccrualResult `json:"successful"`
	Failed       []*InterestAccrualResult `json:"failed"`
	Total        int                      `json:"total"`
	SuccessCount int                      `json:"success_count"`
	FailureCount int                      `json:"failure_count"`
}

type InterestAccrualResult struct {
	AccrualID string  `json:"accrual_id"`
	UserID    string  `json:"user_id"`
	AccountID string  `json:"account_id"`
	Period    string  `json:"period"`
	Interest  float64 `json:"interest"`
	Error     string  `json:"error,omitempty"`
}
