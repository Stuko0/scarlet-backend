package wildfire

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WildfireNRT struct{
	ID primitive.ObjectID `bson:"_id,omitempty"`
	FireID string `bson:"fire_id"`
	Source string `bson:"source"`
	ExternalID string `bson:"external_id"`
	Name string `bson:"name"`
	Status string `bson:"status"`
	DiscoveredTime time.Time `bson:"discovery_time"`
	LastUpdated time.Time `bson:"last_updated"`
	Location GeoPoint `bson:"location"`
	Area float64 `bson:"area"`
	AreaUnit string `bson:"area_unit"`
	Weather WeatherInfo `bson:"weather"`
	MetaData map[string]any `bson:"metadata"`
}

type GeoPoint struct{
	Type string `bson:"type"`
	Coordinates []float64 `bson:"coordinates"`
}

type WeatherInfo struct{
	Temperature float64 `bson:"temperature"`
	WindSpeed float64 `bson:"wind_speed"`
	WindDirection float64 `bson:"wind_direction"`
	Humidity float64 `bson:"humidity"`
	LastUpdated time.Time `bson:"last_updated"`
}