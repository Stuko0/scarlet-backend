package entities

type RTFire  struct{
	Id int64	`json:"id"`
	Latitude   float64 `json:"lat"`
	Longitude  float64 `json:"lng"`
	DetectedAt    string  `json:"detectedAt"`
	Confidence    string  `json:"confidence"`
	FRP        float64 `json:"frp"`
	FWI       float64 `json:"fwi"`
	FireType       string `json:"fireType"`
}