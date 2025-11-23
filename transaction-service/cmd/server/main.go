package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"transaction-service/internal/client"
	"transaction-service/internal/handler"
	"transaction-service/internal/repository/postgres"
	"transaction-service/internal/repository/redis"
	"transaction-service/internal/service"
	"transaction-service/pkg/config"
	"transaction-service/pkg/security"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Инициализация конфигурации
	if err := config.InitConfig(); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	cfg := config.GetConfig()

	// Инициализация базы данных PostgreSQL
	db, err := initPostgres(cfg.Database.PostgresURL)
	if err != nil {
		log.Fatalf("Failed to initialize PostgreSQL: %v", err)
	}
	defer db.Close()

	// Инициализация Redis
	redisClient, err := initRedis(cfg.Database.RedisURL)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	defer redisClient.Close()

	// Инициализация OAuth2 клиента
	oauth2Client := security.NewOAuth2Client(
		cfg.Security.OAuth2.ClientID,
		cfg.Security.OAuth2.ClientSecret,
		cfg.Security.OAuth2.TokenURL,
		[]string{"read", "write"},
	)

	// Инициализация JWT менеджера
	jwtManager := security.NewJWTManager("your-secret-key", 24*time.Hour)

	// Инициализация репозиториев
	transactionRepo := postgres.NewTransactionRepository(db)
	interestRepo := postgres.NewInterestRepository(db)
	notificationRepo := postgres.NewNotificationRepository(db)

	cacheTTL, err := time.ParseDuration(cfg.Cache.TTL)
	if err != nil {
		cacheTTL = 5 * time.Minute
	}
	cacheRepo := redis.NewCacheRepository(redisClient, cacheTTL)

	// Инициализация клиентов
	authClient := client.NewAuthClient(cfg.Services.Auth, oauth2Client)
	depositClient := client.NewDepositClient(cfg.Services.Deposit)

	// Инициализация сервисов уведомлений
	emailService := service.NewEmailService()
	smsService := service.NewSMSService()
	pushService := service.NewPushNotificationService()

	// Инициализация сервисов
	transactionService := service.NewTransactionService(
		transactionRepo,
		cacheRepo,
		authClient,
		depositClient,
	)

	interestService := service.NewInterestService(
		interestRepo,
		transactionRepo,
		cacheRepo,
		depositClient,
	)

	notificationService := service.NewNotificationService(
		notificationRepo,
		cacheRepo,
		emailService,
		smsService,
		pushService,
	)

	// Инициализация обработчиков
	transactionHandler := handler.NewTransactionHandler(transactionService)
	interestHandler := handler.NewInterestHandler(interestService)
	notificationHandler := handler.NewNotificationHandler(notificationService)

	// Настройка маршрутизатора Gin
	router := setupRouter(transactionHandler, interestHandler, notificationHandler, jwtManager)

	// Регистрация в Eureka
	if err := config.RegisterWithEureka(cfg.Eureka.URL, cfg); err != nil {
		log.Printf("Warning: Failed to register with Eureka: %v", err)
	} else {
		log.Printf("Successfully registered with Eureka")
	}

	// Запуск HTTP сервера
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запуск сервера в горутине
	go func() {
		log.Printf("Starting server on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}

func initPostgres(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Проверка соединения
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Настройка пула соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Successfully connected to PostgreSQL")

	// Создание таблиц (в продакшене лучше использовать миграции)
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

func initRedis(redisURL string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		// Если URL не парсится, используем дефолтные настройки
		opts = &redis.Options{
			Addr:     redisURL,
			Password: "", // no password set
			DB:       0,  // use default DB
		}
	}

	client := redis.NewClient(opts)

	// Проверка соединения
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	log.Println("Successfully connected to Redis")
	return client, nil
}

func createTables(db *sql.DB) error {
	// Создание таблицы транзакций
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS transactions (
			id VARCHAR(36) PRIMARY KEY,
			user_id VARCHAR(36) NOT NULL,
			account_id VARCHAR(36) NOT NULL,
			amount DECIMAL(15,2) NOT NULL,
			currency VARCHAR(3) NOT NULL,
			type VARCHAR(20) NOT NULL,
			status VARCHAR(20) NOT NULL,
			description TEXT,
			reference VARCHAR(100),
			metadata JSONB,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create transactions table: %w", err)
	}

	// Создание таблицы начислений процентов
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS interest_accruals (
			id VARCHAR(36) PRIMARY KEY,
			user_id VARCHAR(36) NOT NULL,
			account_id VARCHAR(36) NOT NULL,
			period VARCHAR(7) NOT NULL,
			principal DECIMAL(15,2) NOT NULL,
			interest DECIMAL(15,2) NOT NULL,
			rate DECIMAL(5,2) NOT NULL,
			status VARCHAR(20) NOT NULL,
			applied_at TIMESTAMP,
			created_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create interest_accruals table: %w", err)
	}

	// Создание таблицы процентных ставок
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS interest_rates (
			id VARCHAR(36) PRIMARY KEY,
			account_type VARCHAR(50) NOT NULL,
			rate DECIMAL(5,2) NOT NULL,
			effective_from TIMESTAMP NOT NULL,
			effective_to TIMESTAMP,
			created_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create interest_rates table: %w", err)
	}

	// Создание таблицы уведомлений
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS notifications (
			id VARCHAR(36) PRIMARY KEY,
			user_id VARCHAR(36) NOT NULL,
			type VARCHAR(20) NOT NULL,
			title VARCHAR(255) NOT NULL,
			message TEXT NOT NULL,
			data JSONB,
			status VARCHAR(20) NOT NULL,
			channel VARCHAR(20) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			sent_at TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create notifications table: %w", err)
	}

	// Создание индексов
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
		CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at);
		CREATE INDEX IF NOT EXISTS idx_interest_accruals_user_id ON interest_accruals(user_id);
		CREATE INDEX IF NOT EXISTS idx_interest_accruals_period ON interest_accruals(period);
		CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
		CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at);
	`)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	log.Println("Database tables created/verified successfully")
	return nil
}

func setupRouter(
	transactionHandler *handler.TransactionHandler,
	interestHandler *handler.InterestHandler,
	notificationHandler *handler.NotificationHandler,
	jwtManager *security.JWTManager,
) *gin.Engine {
	router := gin.Default()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(security.JWTAuthMiddleware(jwtManager))

	// Группа маршрутов для транзакций
	transactions := router.Group("/api/transactions")
	{
		transactions.POST("", transactionHandler.CreateTransaction)
		transactions.POST("/batch", transactionHandler.BatchCreateTransactions)
		transactions.GET("/:id", transactionHandler.GetTransaction)
	}

	// Группа маршрутов для пользователей
	users := router.Group("/api/users")
	{
		users.GET("/:user_id/transactions", transactionHandler.GetUserTransactions)
		users.GET("/:user_id/transactions/stats", transactionHandler.GetTransactionStats)
		users.GET("/:user_id/interest/history", interestHandler.GetInterestHistory)
		users.GET("/:user_id/notifications", notificationHandler.GetUserNotifications)
		users.GET("/:user_id/notifications/stats", notificationHandler.GetNotificationStats)
	}

	// Группа маршрутов для процентов
	interest := router.Group("/api/interest")
	{
		interest.POST("/calculate", interestHandler.CalculateInterest)
		interest.POST("/apply", interestHandler.ApplyInterest)
		interest.POST("/process-pending", security.RequireRoles("ADMIN"), interestHandler.ProcessPendingInterest)
	}

	// Группа маршрутов для уведомлений
	notifications := router.Group("/api/notifications")
	{
		notifications.POST("", notificationHandler.CreateNotification)
		notifications.POST("/transaction", notificationHandler.SendTransactionNotification)
		notifications.PUT("/:id/mark-sent", notificationHandler.MarkAsSent)
	}

	// Health check
	router.GET("/health", transactionHandler.HealthCheck)

	// API documentation
	router.GET("/swagger/*any", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Swagger documentation will be available here",
		})
	})

	return router
}
