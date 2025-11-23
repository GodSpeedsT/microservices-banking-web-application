package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type DepositClient struct {
	baseURL    string
	httpClient *http.Client
}

type AccountBalance struct {
	AccountID string  `json:"account_id"`
	UserID    string  `json:"user_id"`
	Balance   float64 `json:"balance"`
	Currency  string  `json:"currency"`
	Available float64 `json:"available_balance"`
	Locked    float64 `json:"locked_balance"`
	UpdatedAt string  `json:"updated_at"`
}

type TransactionRequest struct {
	FromAccountID string  `json:"from_account_id,omitempty"`
	ToAccountID   string  `json:"to_account_id,omitempty"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Type          string  `json:"type"`
	Description   string  `json:"description"`
	Reference     string  `json:"reference"`
}

type TransactionResponse struct {
	TransactionID string  `json:"transaction_id"`
	Status        string  `json:"status"`
	NewBalance    float64 `json:"new_balance"`
	Message       string  `json:"message"`
}

func NewDepositClient(baseURL string) *DepositClient {
	return &DepositClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *DepositClient) GetBalance(ctx context.Context, accountID string, accessToken string) (*AccountBalance, error) {
	url := fmt.Sprintf("%s/api/accounts/%s/balance", c.baseURL, accountID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("deposit service returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var balance AccountBalance
	if err := json.Unmarshal(body, &balance); err != nil {
		return nil, fmt.Errorf("failed to parse balance: %w", err)
	}

	return &balance, nil
}

func (c *DepositClient) ProcessTransaction(ctx context.Context, req *TransactionRequest, accessToken string) (*TransactionResponse, error) {
	url := fmt.Sprintf("%s/api/transactions", c.baseURL)

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+accessToken)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("deposit service returned status %d: %s", resp.StatusCode, string(body))
	}

	var transactionResp TransactionResponse
	if err := json.Unmarshal(body, &transactionResp); err != nil {
		return nil, fmt.Errorf("failed to parse transaction response: %w", err)
	}

	return &transactionResp, nil
}

func (c *DepositClient) HealthCheck(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/health", nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status: %d", resp.StatusCode)
	}

	return nil
}
