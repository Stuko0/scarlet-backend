package model

type RTFire struct {
	DocID      string  `json:"-"`
	Id         int64   `json:"id"`
	Latitude   float64 `json:"lat"`
	Longitude  float64 `json:"lng"`
	Location   string  `json:"location"`
	DetectedAt string  `json:"detectedAt"`
	Confidence string  `json:"confidence"`
	FRP        float64 `json:"frp"`
	FWI        float64 `json:"fwi"`
	FireType   string  `json:"fireType"`
	Active     string  `json:"active"`
}
