package database

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisPubSub struct{
	client *redis.Client
}

func NewRedisClient() *redis.Client{
	return redis.NewClient(&redis.Options{
		Addr: "localhost:8000",
		Password: "",
		DB: 0,
	})
}

func NewRedisPubSub(client *redis.Client)*RedisPubSub{
	return &RedisPubSub{client: client}
}

func (r *RedisPubSub) Publish(ctx context.Context, channel string, message interface{})error{
	return r.client.Publish(ctx, channel, message).Err()
}

func (r *RedisPubSub) Suscribe(ctx context.Context, channels ...string)*redis.PubSub{
	return r.client.Subscribe(ctx, channels...)
}