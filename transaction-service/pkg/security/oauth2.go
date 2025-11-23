package security

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"io"
	"net/http"
	"strings"
)

type OAuth2Client struct {
	config *clientcredentials.Config
	client *http.Client
}

type TokenInfo struct {
	Active   bool     `json:"active"`
	UserID   string   `json:"user_id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
	ClientID string   `json:"client_id"`
	Exp      int64    `json:"exp"`
	Scope    string   `json:"scope"`
}

func NewOAuth2Client(clientID, clientSecret, tokenURL string, scopes []string) *OAuth2Client {
	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
		Scopes:       scopes,
	}

	return &OAuth2Client{
		config: config,
		client: config.Client(context.Background()),
	}
}

func (c *OAuth2Client) ValidateToken(ctx context.Context, token string) (*TokenInfo, error) {
	// В реальном приложении здесь будет вызов к auth-service для проверки токена
	// Для примера делаем упрощенную реализацию

	req, err := http.NewRequestWithContext(ctx, "GET",
		strings.Replace(c.config.TokenURL, "token", "check_token", 1), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create validation request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token validation failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var tokenInfo TokenInfo
	if err := json.Unmarshal(body, &tokenInfo); err != nil {
		return nil, fmt.Errorf("failed to parse token info: %w", err)
	}

	if !tokenInfo.Active {
		return nil, errors.New("token is not active")
	}

	return &tokenInfo, nil
}

func (c *OAuth2Client) GetClientCredentialsToken(ctx context.Context) (*oauth2.Token, error) {
	return c.config.Token(ctx)
}

// Middleware для OAuth2 аутентификации
func OAuth2AuthMiddleware(oauthClient *OAuth2Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(401, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]
		tokenInfo, err := oauthClient.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token", "message": err.Error()})
			c.Abort()
			return
		}

		// Сохраняем информацию о пользователе в контексте
		c.Set("userInfo", tokenInfo)
		c.Next()
	}
}
