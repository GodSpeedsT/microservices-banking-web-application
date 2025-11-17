// cmd/server/main.go
package main

import (
	"log"
	"transaction-service/internal/handler"
	"transaction-service/internal/service"
	"transaction-service/pkg/config"

	"github.com/gin-gonic/gin"
)

func main() {
	// Загрузка конфигурации
	cfg := config.LoadConfig()

	// Инициализация сервисов
	transactionService := service.NewTransactionService()
	interestService := service.NewInterestService()
	notificationService := service.NewNotificationService()

	// Инициализация обработчиков
	transactionHandler := handler.NewTransactionHandler(transactionService)
	interestHandler := handler.NewInterestHandler(interestService)
	notificationHandler := handler.NewNotificationHandler(notificationService)

	// HTTP роутер
	router := gin.Default()

	// Middleware для OAuth2 аутентификации
	router.Use(OAuth2Middleware())

	// Routes
	api := router.Group("/api")
	{
		api.POST("/transactions/process", transactionHandler.ProcessTransaction)
		api.POST("/transactions/calculate-interest", interestHandler.CalculateInterest)
		api.POST("/notifications/send", notificationHandler.SendNotification)
		api.GET("/health", healthCheck)
	}

	// Регистрация в Eureka
	config.RegisterWithEureka(cfg)

	log.Printf("Transaction Service starting on port %s", cfg.Server.Port)
	router.Run(":" + cfg.Server.Port)
}

func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "UP"})
}

func OAuth2Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if !validateToken(token) {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
