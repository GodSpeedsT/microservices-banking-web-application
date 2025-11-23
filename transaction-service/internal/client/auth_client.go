package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"transaction-service/pkg/security"
)

type AuthClient struct {
	baseURL      string
	httpClient   *http.Client
	oauth2Client *security.OAuth2Client
}

type UserInfo struct {
	UserID   string   `json:"user_id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
	Enabled  bool     `json:"enabled"`
	Active   bool     `json:"active"`
}

func NewAuthClient(baseURL string, oauth2Client *security.OAuth2Client) *AuthClient {
	return &AuthClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		oauth2Client: oauth2Client,
	}
}

func (c *AuthClient) ValidateToken(ctx context.Context, token string) (*UserInfo, error) {
	// Используем OAuth2 клиент для проверки токена
	tokenInfo, err := c.oauth2Client.ValidateToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	return &UserInfo{
		UserID:   tokenInfo.UserID,
		Username: tokenInfo.Username,
		Email:    tokenInfo.Email,
		Roles:    tokenInfo.Roles,
		Enabled:  true,
		Active:   tokenInfo.Active,
	}, nil
}

func (c *AuthClient) GetUserByID(ctx context.Context, userID string) (*UserInfo, error) {
	// Получаем client credentials token для межсервисного взаимодействия
	token, err := c.oauth2Client.GetClientCredentialsToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get client credentials token: %w", err)
	}

	url := fmt.Sprintf("%s/api/users/%s", c.baseURL, userID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth service returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var userInfo UserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	return &userInfo, nil
}

func (c *AuthClient) HasRole(ctx context.Context, token string, requiredRole string) (bool, error) {
	userInfo, err := c.ValidateToken(ctx, token)
	if err != nil {
		return false, err
	}

	for _, role := range userInfo.Roles {
		if role == requiredRole {
			return true, nil
		}
	}

	return false, nil
}

func (c *AuthClient) HealthCheck(ctx context.Context) error {
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
