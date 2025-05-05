package wildfire

import (
	"context"
	"fmt"
	"log"
	"time"

	"connectrpc.com/connect"
	nrtv1 "github.com/Stuko0/scarlet-backend/gen/proto/wildfire/v1"
	"github.com/Stuko0/scarlet-backend/internal/domain/wildfire/scrapers"
	"go.mongodb.org/mongo-driver/mongo"
)

type WildfireNRTService struct{
	repo WildfireNRTRepository
	scrapers []scrapers.Scraper
	weather WeatherProvider
}

type WeatherProvider interface {
    GetWeatherData(ctx context.Context, lat, lng float64) (*scrapers.WeatherData, error)
}

func NewWildfireNRTService(repo WildfireNRTRepository, scrapers []scrapers.Scraper, weather WeatherProvider) *WildfireNRTService{
	return &WildfireNRTService{
		repo: repo,
		scrapers: scrapers,
		weather: weather,
	}
}

func ProtoToWildfireNRT(p *nrtv1.WildfireNRT)(*WildfireNRT, error){
	discoveryTime, err:=time.Parse(time.RFC3339, p.DiscoveryTime)
	if err!=nil{
		return nil, err
	}

	lastUpdated, err:= time.Parse(time.RFC3339, p.LastUpdated)
	if err!=nil{
		return nil, err
	}
	weatherLastUpdated, err:=time.Parse(time.RFC3339, p.Weather.LastUpdated)
	if err!=nil{
		return nil, err
	}

	return &WildfireNRT{
		FireID: p.FireId,
		Source: p.Source,
		ExternalID: p.ExternalId,
		Name: p.Name,
		Status: p.Status,
		DiscoveredTime: discoveryTime,
		LastUpdated: lastUpdated,
		Location: GeoPoint{
			Type: p.Location.Type,
			Coordinates: p.Location.Coordinates,
		},
		Area: p.Area,
		AreaUnit: p.AreaUnit,
		Weather: WeatherInfo{
			Temperature: p.Weather.Temperature,
			WindSpeed: p.Weather.WindSpeed,
			WindDirection: p.Weather.WindDirection,
			Humidity: p.Weather.Humidity,
			LastUpdated: weatherLastUpdated,
		},
		MetaData: p.Metadata,
	}, nil
}

func WildfireNRTToProto(w *WildfireNRT)*nrtv1.WildfireNRT{
	return &nrtv1.WildfireNRT{
		FireId: w.FireID,
		Source: w.Source,
		ExternalId: w.ExternalID,
		Name: w.Name,
		Status: w.Status,
		DiscoveryTime: w.DiscoveredTime.Format(time.RFC3339),
		LastUpdated: w.LastUpdated.Format(time.RFC3339),
		Location: &nrtv1.GeoPoint{
			Type: w.Location.Type,
			Coordinates: w.Location.Coordinates,
		},
		Area: w.Area,
		AreaUnit: w.AreaUnit,
		Weather: &nrtv1.WeatherData{
			Temperature: w.Weather.Temperature,
			WindSpeed: w.Weather.WindSpeed,
			WindDirection: w.Weather.WindDirection,
			Humidity: w.Weather.Humidity,
			LastUpdated: w.Weather.LastUpdated.Format(time.RFC3339),
		},
		Metadata: w.MetaData,
	}
}

func (s *WildfireNRTService) UpsertFromDetection(ctx context.Context, req *connect.Request[nrtv1.UpsertRequest])(*connect.Response[nrtv1.UpsertResponse], error){
	wildfire, err:= ProtoToWildfireNRT(req.Msg.Wildfire)
	if err!=nil{
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	err= s.repo.UpsertFromDetection(ctx, wildfire)
	if err!=nil{
		return nil, connect.NewError(connect.CodeInternal,err)
	}

	return connect.NewResponse(&nrtv1.UpsertResponse{
		DocId: wildfire.ID.Hex(),
	}), nil
}

func (s *WildfireNRTService) GetActiveFires(ctx context.Context, req *connect.Request[nrtv1.GetActiveFiresRequest])(*connect.Response[nrtv1.GetActiveFiresResponse], error){
	fires, err:= s.repo.GetActiveFires(ctx)
	if err!=nil{
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	protoFires:=make([]*nrtv1.WildfireNRT, len(fires))
	for i, fire := range fires{
		protoFires[i]=WildfireNRTToProto(fire)
	}

	return connect.NewResponse(&nrtv1.GetActiveFiresResponse{
		Wildfires: protoFires,
	}), nil
}

func (s*WildfireNRTService) GetNearbyFires(ctx context.Context, req *connect.Request[nrtv1.GetNearbyFiresRequest])(*connect.Response[nrtv1.GetNearbyFiresResponse], error){
	if len(req.Msg.Location.Coordinates)!=2{
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("location must have 2 latitude and longitude"))
	}
	latitude := req.Msg.Location.Coordinates[0]
	longitude := req.Msg.Location.Coordinates[1]
	fires, err:= s.repo.GetNearbyFires(ctx, latitude, longitude, int(req.Msg.RadiusKm))
	if err!=nil{
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	protoFires:= make([]*nrtv1.WildfireNRT, len(fires))
	for i, fire:=range fires{
		protoFires[i]=WildfireNRTToProto(fire)
	}

	return connect.NewResponse(&nrtv1.GetNearbyFiresResponse{
		Wildfires: protoFires,
	}), nil
}

func (s *WildfireNRTService) processScrape(ctx context.Context, scraper scrapers.Scraper){
	detections, err:= scraper.Scrape(ctx)
	if err!=nil{
		log.Printf("scraper %s failed: %v", scraper.Name(), err)
		return
	}

	bulkCtx, cancel:=context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fires:=make([]*WildfireNRT,0,len(detections))
	weatherFailures:=0

	for _, detection:= range detections{
		weather, err:= s.weather.GetWeatherData(bulkCtx, detection.Latitude, detection.Longitude)
		if err!=nil{
			weatherFailures++
			log.Printf("failed to get weather for %s: %v", detection.ExternalID, err)
			continue
		}
		fires= append(fires, &WildfireNRT{
			FireID: detection.ExternalID,
			Source: detection.Source,
			ExternalID: detection.ExternalID,
			DiscoveredTime: detection.DiscoveredAt,
			LastUpdated: time.Now(),
			Location: GeoPoint{
				Type: "Point",
				Coordinates: []float64{detection.Longitude, detection.Latitude},
			},
			Weather: WeatherInfo{
				Temperature: weather.Temperature,
				WindSpeed: weather.WindSpeed,
				WindDirection: weather.WindDirection,
				Humidity: weather.Humidity,
				LastUpdated: time.Now(),
			},
			MetaData: detection.AdditionalData,
		})
	}

	const maxBatchSize=100
	totalFires:=len(fires)
	for i:=0;i<len(fires);i+=maxBatchSize{
		end:=i+maxBatchSize
		if end>totalFires{end = totalFires}
		batch:=fires[i:end]

		if err:= s.repo.BulkUpsertFromDetections(bulkCtx, batch); err!=nil{
			if mongo.IsDuplicateKeyError(err) {
				log.Printf("duplicate fires in batch [%d-%d] (non-fatal)", i, end-1)
			}else{log.Printf("failed to upsert batch [%d-%d]: %v",i,end-1,err)}
		}else{
			log.Printf("successfully upsert batch [%d-%d](%d fires)",i,end-1,len(batch))
		}
	}
	if weatherFailures>0{log.Printf("completed with %d weather API failures", weatherFailures)}
}

func (s *WildfireNRTService) runScraper(ctx context.Context, scraper scrapers.Scraper){
	ticker:= time.NewTicker(scraper.Interval())
	defer ticker.Stop()

	for{
		select{
		case<-ticker.C:
			s.processScrape(ctx, scraper)
		case<-ctx.Done():
			return
		}
	}
}

func (s *WildfireNRTService) TriggerImmediateScrape(ctx context.Context){
	for _, scraper:=range s.scrapers{s.processScrape(ctx,scraper)}
}

func (s *WildfireNRTService) RunScrapers(ctx context.Context){
	for _, scraper := range s.scrapers{
		go s.runScraper(ctx, scraper)
	}
}