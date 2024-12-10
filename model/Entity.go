package model

type Entity struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
	Location  string  `json:"location"`
	Social    string  `json:"social"`
	Image     string  `json:"image_url"`
	Active    string  `json:"active"`
}
