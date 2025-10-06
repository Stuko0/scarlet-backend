package wildfire

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Stuko0/scarlet-backend/internal/database"
)

type CachedRepository struct {
	mongoRepo WildfireNRTRepository
	redis *database.RedisClient
}

func NewCachedRepository(mongoRepo WildfireNRTRepository, redis *database.RedisClient) *CachedRepository {
	return &CachedRepository{
		mongoRepo: mongoRepo,
		redis:     redis,
	}
}


func (r *CachedRepository) GetActiveFires(ctx context.Context) ([]*WildfireNRT, error) {
    // 1. Intenta obtener de Redis
    var cached []*WildfireNRT
    if err := r.redis.GetJSON(ctx, "wildfires:active", &cached); err == nil {
        log.Println("Cache hit - returning from Redis")
        return cached, nil
    }
    
    log.Println("Cache miss - querying MongoDB")
    
    // 2. Si Redis falla, ve a MongoDB
    fires, err := r.mongoRepo.GetActiveFires(ctx)
    if err != nil {
        return nil, err
    }
    
    // 3. Guarda en Redis (solo si hay datos)
    if len(fires) > 0 {
        if err := r.redis.SetJSON(ctx, "wildfires:active", fires, 3*time.Hour); err != nil {
            log.Printf("Failed to cache fires: %v", err)
        } else {
            log.Printf("Cached %d fires in Redis", len(fires))
        }
    }
    
    return fires, nil
}

// Implementa los demás métodos del interfaz WildfireNRTRepository...
func (r *CachedRepository) UpsertFromDetection(ctx context.Context, detection *WildfireNRT) error {
	err := r.mongoRepo.UpsertFromDetection(ctx, detection)
	if err == nil {
		r.redis.PublishJSON(ctx, "wildfires:updates", detection)
	}
	return err
}

func (r *CachedRepository) GetNearbyFires(ctx context.Context, lat, lng float64, radiusKm int) ([]*WildfireNRT, error) {
	// Intenta obtener de Redis primero
	var cachedFires []*WildfireNRT
	cacheKey := "wildfires:nearby:" + fmt.Sprintf("%f:%f:%d", lat, lng, radiusKm)
	err := r.redis.GetJSON(ctx, cacheKey, &cachedFires)
	if err == nil {
		return cachedFires, nil
	}

	// Si falla, va a MongoDB
	fires, err := r.mongoRepo.GetNearbyFires(ctx, lat, lng, radiusKm)
	if err != nil {
		return nil, err
	}

	// Almacena en caché
	if err := r.redis.SetJSON(ctx, cacheKey, fires, 15*time.Minute); err != nil {
		log.Printf("Failed to cache nearby fires: %v", err)
	}

	return fires, nil
}

func (r *CachedRepository) BulkUpsertFromDetections(ctx context.Context, detections []*WildfireNRT) error {
    // 1. Actualiza MongoDB primero
    if err := r.mongoRepo.BulkUpsertFromDetections(ctx, detections); err != nil {
        return fmt.Errorf("mongo bulk upsert failed: %w", err)
    }

    // 2. Invalida la caché de manera eficiente
    if err := r.redis.Delete(ctx, "wildfires:active"); err != nil {
        log.Printf("Cache invalidation warning: %v", err)
        // No retornes error aquí para no fallar toda la operación
    }

    // 3. Publica actualizaciones en segundo plano
    go func() {
        for _, detection := range detections {
            if err := r.redis.PublishJSON(context.Background(), "wildfires:updates", detection); err != nil {
                log.Printf("Failed to publish update for %s: %v", detection.FireID, err)
            }
        }
    }()
    
    return nil
}