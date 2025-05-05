package scrapers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type OpenMeteoScraper struct{
	baseURL string
}

func NewOpenMeteoScraper() *OpenMeteoScraper{
	return &OpenMeteoScraper{
		baseURL: "https://api.open-meteo.com/v1",
	}
}

type WeatherData struct{
	Temperature float64
	WindSpeed float64
	WindDirection float64
	Humidity float64
	ForecastHours int
}

func (s *OpenMeteoScraper) GetWeatherData(ctx context.Context, lat, lng float64) (*WeatherData, error){
	url := fmt.Sprintf(
		"%s/forecast?latitude=%f&longitude=%f&current_weather=true&hourly=temperature_2m,relativehumidity_2m,windspeed_10m",
		s.baseURL, lat, lng,
	)

	resp, err:=http.Get(url)
	if err!=nil{
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	var result struct{
		Current struct{
			Temperature float64 `json:"temperature"`
			WindSpeed float64 `json:"windspeed"`
			WindDir float64 `json:"winddirection"`
		} `json:"current_weather"`
		Hourly struct{
			Time []string `json:"time"`
			Temperature []float64 `json:"temperature_2m"`
			Humidity []float64 `json:"relativehumidity_2m"`
			WindSpeed []float64 `json:"windspeed_10m"`
		} `json:"hourly"`
	}
	if err:= json.NewDecoder(resp.Body).Decode(&result); err!=nil{
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &WeatherData{
		Temperature: result.Current.Temperature,
		WindSpeed: result.Current.WindSpeed,
		WindDirection: result.Current.WindDir,
		Humidity: result.Hourly.Humidity[0],
		ForecastHours: len(result.Hourly.Time),
	}, nil
}