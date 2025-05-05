package scrapers

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type FIRMSScaper struct{
	apiKey string
	country string
}

func NewFIRMSScraper(apiKey, country string) *FIRMSScaper{
	return &FIRMSScaper{
		apiKey: apiKey,
		country: country,
	}
}
func (s *FIRMSScaper) Name() string { return "nasa_firms"}
func (s *FIRMSScaper) Interval() time.Duration {return 1* time.Hour}
func (s *FIRMSScaper) Scrape(ctx context.Context) ([]*WildfireDetection, error){
	url := fmt.Sprintf("https://firms.modaps.eosdis.nasa.gov/api/country/csv/%s/VIIRS_SNPP_NRT/%s/1", s.apiKey, s.country)

	resp, err:= http.Get(url)
	if err!=nil{return nil, fmt.Errorf("failed to make request: %w", err)}

	
	defer resp.Body.Close()

	if resp.StatusCode!=200{return nil, fmt.Errorf("API request failed: HTTP %d", resp.StatusCode)}

	reader:= csv.NewReader(resp.Body)
	reader.Comma=','
	reader.Comment='#'
	if _,err:=reader.Read(); err!=nil{
		return nil, fmt.Errorf("failed to read header: %w", err)
	}
	var detections []*WildfireDetection

	for{
		record, err:= reader.Read()
		if err == io.EOF{break}
		if err!= nil{
			return nil, fmt.Errorf("failed to read record: %w", err)
		}
		lat, _ := strconv.ParseFloat(record[1],64)
		lng, _ := strconv.ParseFloat(record[2],64)
		confidence, _ := strconv.ParseFloat(record[10],64)
		dateTimeStr:=fmt.Sprintf("%s-%s", record[6], record[7])
		date, _ := time.Parse("2006-01-02-1504", dateTimeStr)

		detections=append(detections, &WildfireDetection{
			Source: s.Name(),
			ExternalID: fmt.Sprintf("%s-%s-%s-%s", record[6], record[7], record[1], record[2]),
			Latitude: lat,
			Longitude: lng,
			Confidence: confidence,
			DiscoveredAt: date,
			AdditionalData: map[string]interface{}{
				"brightness" : record[3],
				"frp" : record[13],
				"daynight" : record[14],
			},
		})
	}
	return detections, nil
}