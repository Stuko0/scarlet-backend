package scrapers

import (
	"context"
	"time"
)

type Scraper interface{
	Name() string
	Scrape(ctx context.Context)([]*WildfireDetection, error)
	Interval() time.Duration
}

type WildfireDetection struct{
	Source string
	ExternalID string
	Latitude float64
	Longitude float64
	DiscoveredAt time.Time
	Confidence float64
	AdditionalData map[string]interface{}
}