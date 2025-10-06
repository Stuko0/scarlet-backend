package wildfire

import "context"

type WildfireNRTRepository interface {
	UpsertFromDetection(ctx context.Context, detection *WildfireNRT) error
	GetActiveFires(ctx context.Context) ([]*WildfireNRT, error)
	GetNearbyFires(ctx context.Context, lat, lng float64, radiusKm int) ([]*WildfireNRT, error)
	BulkUpsertFromDetections(ctx context.Context, detections []*WildfireNRT) error
}
type CacheRepository interface {
	CacheActiveFires(ctx context.Context, fires []*WildfireNRT) error
	GetCachedActiveFires(ctx context.Context) ([]*WildfireNRT, error)
	GetCachedNearbyFires(ctx context.Context, lat, lng float64, radiusKm int) ([]*WildfireNRT, error)
	ClearCache(ctx context.Context) error
	PublishUpdate(ctx context.Context, update interface{}) error
}