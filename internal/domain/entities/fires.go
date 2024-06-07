package entities

type Fire  struct{
	Id int64	`json:"id"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	AcqDate    string  `json:"acq_date"`
	AcqTime    string  `json:"acq_time"`
	Brightness float64 `json:"brightness"`
	Version    string  `json:"version"`
	BrightT31  float64 `json:"bright_t31"`
	Confidence string  `json:"confidence"`
	DayNight   string  `json:"daynight"`
	Instrument string  `json:"instrument"`
	FRP        float64 `json:"frp"`
	Scan       float64 `json:"scan"`
	Track      float64 `json:"track"`
	Satellite  string  `json:"satellite"`
}