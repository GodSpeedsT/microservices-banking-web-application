package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheRepository struct {
	client *redis.Client
	ttl    time.Duration
}

func NewCacheRepository(client *redis.Client, ttl time.Duration) *CacheRepository {
	return &CacheRepository{
		client: client,
		ttl:    ttl,
	}
}

func (r *CacheRepository) Set(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	err = r.client.Set(ctx, key, data, r.ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

func (r *CacheRepository) Get(ctx context.Context, key string, dest interface{}) (bool, error) {
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, fmt.Errorf("failed to get cache: %w", err)
	}

	if err := json.Unmarshal([]byte(data), dest); err != nil {
		return false, fmt.Errorf("failed to unmarshal cached data: %w", err)
	}

	return true, nil
}

func (r *CacheRepository) Delete(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete cache: %w", err)
	}
	return nil
}

func (r *CacheRepository) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check cache existence: %w", err)
	}
	return count > 0, nil
}

func (r *CacheRepository) SetWithCustomTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	err = r.client.Set(ctx, key, data, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set cache with custom TTL: %w", err)
	}

	return nil
}

func (r *CacheRepository) HSet(ctx context.Context, key string, values map[string]interface{}) error {
	err := r.client.HSet(ctx, key, values).Err()
	if err != nil {
		return fmt.Errorf("failed to hset cache: %w", err)
	}

	err = r.client.Expire(ctx, key, r.ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set expiration for hset: %w", err)
	}

	return nil
}

func (r *CacheRepository) HGet(ctx context.Context, key, field string, dest interface{}) (bool, error) {
	data, err := r.client.HGet(ctx, key, field).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, fmt.Errorf("failed to hget cache: %w", err)
	}

	if err := json.Unmarshal([]byte(data), dest); err != nil {
		return false, fmt.Errorf("failed to unmarshal cached hget data: %w", err)
	}

	return true, nil
}

func (r *CacheRepository) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	result, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to hgetall cache: %w", err)
	}

	if len(result) == 0 {
		return nil, redis.Nil
	}

	return result, nil
}

func (r *CacheRepository) DeletePattern(ctx context.Context, pattern string) error {
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get keys by pattern: %w", err)
	}

	if len(keys) > 0 {
		err = r.client.Del(ctx, keys...).Err()
		if err != nil {
			return fmt.Errorf("failed to delete keys by pattern: %w", err)
		}
	}

	return nil
}

// Специфичные методы для transaction service

func (r *CacheRepository) CacheUserTransactions(ctx context.Context, userID string, transactions interface{}) error {
	key := fmt.Sprintf("user:%s:transactions", userID)
	return r.Set(ctx, key, transactions)
}

func (r *CacheRepository) GetCachedUserTransactions(ctx context.Context, userID string, dest interface{}) (bool, error) {
	key := fmt.Sprintf("user:%s:transactions", userID)
	return r.Get(ctx, key, dest)
}

func (r *CacheRepository) CacheTransaction(ctx context.Context, transactionID string, transaction interface{}) error {
	key := fmt.Sprintf("transaction:%s", transactionID)
	return r.Set(ctx, key, transaction)
}

func (r *CacheRepository) GetCachedTransaction(ctx context.Context, transactionID string, dest interface{}) (bool, error) {
	key := fmt.Sprintf("transaction:%s", transactionID)
	return r.Get(ctx, key, dest)
}

func (r *CacheRepository) InvalidateUserTransactions(ctx context.Context, userID string) error {
	pattern := fmt.Sprintf("user:%s:transactions*", userID)
	return r.DeletePattern(ctx, pattern)
}

func (r *CacheRepository) InvalidateTransaction(ctx context.Context, transactionID string) error {
	key := fmt.Sprintf("transaction:%s", transactionID)
	return r.Delete(ctx, key)
}
