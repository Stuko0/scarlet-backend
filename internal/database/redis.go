package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct{
	client *redis.Client
	ttl time.Duration
}

func NewRedisClient() *RedisClient{
	client := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
        PoolSize: 100,
    })
    
    // Verifica la conexión
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := client.Ping(ctx).Err(); err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }
    
    log.Println("Successfully connected to Redis")
    
    return &RedisClient{
        client: client,
        ttl:    1 * time.Hour,
    }
}

func (r *RedisClient) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
    // Agrega logs para depuración
    log.Printf("Caching data for key: %s", key)
    
    jsonData, err := json.Marshal(value)
    if err != nil {
        log.Printf("Failed to marshal JSON for Redis: %v", err)
        return fmt.Errorf("marshal error: %w", err)
    }
    
    if err := r.client.Set(ctx, key, jsonData, ttl).Err(); err != nil {
        log.Printf("Failed to set key %s in Redis: %v", key, err)
        return fmt.Errorf("redis set error: %w", err)
    }
    
    log.Printf("Successfully cached data for key: %s", key)
    return nil
}

func (r *RedisClient) GetJSON(ctx context.Context, key string, dest interface{}) error {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func (r *RedisClient) PublishJSON(ctx context.Context, channel string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Publish(ctx, channel, data).Err()
}

func (r *RedisClient) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return r.client.Subscribe(ctx, channels...)
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}

func (r *RedisClient) Delete(ctx context.Context, key string) error {
    if err := r.client.Del(ctx, key).Err(); err != nil {
        return fmt.Errorf("failed to delete key %s: %w", key, err)
    }
    return nil
}