module transaction-service

go 1.19

require (
    // HTTP сервер
    github.com/gin-gonic/gin v1.9.1

    // Базы данных
    github.com/lib/pq v1.10.9                    // PostgreSQL
    github.com/go-redis/redis/v8 v8.11.5         // Redis

    // Конфигурация
    github.com/spf13/viper v1.16.0

    // Утилиты
    github.com/sirupsen/logrus v1.9.3           // Логирование
    github.com/google/uuid v1.3.0               // UUID генерация
    golang.org/x/crypto v0.9.0                  // Crypto утилиты
)