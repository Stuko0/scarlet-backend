package wildfire

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	"connectrpc.com/connect"
// 	nrtv1 "github.com/Stuko0/scarlet-backend/gen/proto/wildfire/v1"
// 	"github.com/Stuko0/scarlet-backend/internal/domain/wildfire/scrapers"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/suite"
// )

// type MockWildfireRepository struct {
// 	mock.Mock
// }

// func (m *MockWildfireRepository) UpsertFromDetection(ctx context.Context, detection *WildfireNRT) error {
// 	args := m.Called(ctx, detection)
// 	return args.Error(0)
// }

// func (m *MockWildfireRepository) BulkUpsertFromDetections(ctx context.Context, detections []*WildfireNRT) error {
// 	args := m.Called(ctx, detections)
// 	return args.Error(0)
// }

// func (m *MockWildfireRepository) GetActiveFires(ctx context.Context) ([]*WildfireNRT, error) {
// 	args := m.Called(ctx)
// 	return args.Get(0).([]*WildfireNRT), args.Error(1)
// }

// func (m *MockWildfireRepository) GetNearbyFires(ctx context.Context, lat, lng float64, radiusKm int) ([]*WildfireNRT, error) {
// 	args := m.Called(ctx, lat, lng, radiusKm)
// 	return args.Get(0).([]*WildfireNRT), args.Error(1)
// }

// type MockScraper struct {
// 	mock.Mock
// }

// func (m *MockScraper) Name() string {
// 	return "mock-scraper"
// }

// func (m *MockScraper) Interval() time.Duration {
// 	return 5 * time.Minute
// }

// func (m *MockScraper) Scrape(ctx context.Context) ([]*scrapers.WildfireDetection, error) {
// 	args := m.Called(ctx)
// 	return args.Get(0).([]*scrapers.WildfireDetection), args.Error(1)
// }

// type MockWeatherProvider struct {
// 	mock.Mock
// }

// func (m *MockWeatherProvider) GetWeatherData(ctx context.Context, lat, lng float64) (*scrapers.WeatherData, error) {
// 	args := m.Called(ctx, lat, lng)
// 	return args.Get(0).(*scrapers.WeatherData), args.Error(1)
// }

// type WildfireServiceTestSuite struct {
// 	suite.Suite
// 	ctx        context.Context
// 	service    *WildfireNRTService
// 	mockRepo   *MockWildfireRepository
// 	mockScraper *MockScraper
// 	mockWeather *MockWeatherProvider
// }

// func (s *WildfireServiceTestSuite) SetupTest() {
// 	s.ctx = context.Background()
// 	s.mockRepo = new(MockWildfireRepository)
// 	s.mockScraper = new(MockScraper)
// 	s.mockWeather = new(MockWeatherProvider)
	
// 	s.service = NewWildfireNRTService(
// 		s.mockRepo,
// 		[]scrapers.Scraper{s.mockScraper},
// 		s.mockWeather,
// 	)
// }

// func (s *WildfireServiceTestSuite) TestGetActiveFires_Success() {
// 	testFires := []*WildfireNRT{
// 		{
// 			FireID: "fire-123",
// 			Location: GeoPoint{
// 				Type:        "Point",
// 				Coordinates: []float64{-118.5, 34.0},
// 			},
// 		},
// 	}

// 	// Set expectations
// 	s.mockRepo.On("GetActiveFires", s.ctx).Return(testFires, nil)

// 	// Execute
// 	req := &connect.Request[nrtv1.GetActiveFiresRequest]{}
// 	resp, err := s.service.GetActiveFires(s.ctx, req)

// 	// Verify
// 	s.NoError(err)
// 	s.NotNil(resp)
// 	s.Len(resp.Msg.Wildfires, 1)
// 	s.Equal("fire-123", resp.Msg.Wildfires[0].FireId)
// 	s.mockRepo.AssertExpectations(s.T())
// }

// func (s *WildfireServiceTestSuite) TestGetNearbyFires_Success() {
// 	// Setup test data
// 	testFires := []*WildfireNRT{
// 		{
// 			FireID: "fire-nearby",
// 			Location: GeoPoint{
// 				Type:        "Point",
// 				Coordinates: []float64{-118.4, 34.1},
// 			},
// 		},
// 	}

// 	// Set expectations - note correct parameter order (lat, lng)
// 	s.mockRepo.On("GetNearbyFires", s.ctx, -118.5, 34.0, 10).Return(testFires, nil)

// 	// Execute
// 	req := &connect.Request[nrtv1.GetNearbyFiresRequest]{
// 		Msg: &nrtv1.GetNearbyFiresRequest{
// 			Location: &nrtv1.GeoPoint{
// 				Coordinates: []float64{-118.5, 34.0}, // Note: longitude first
// 			},
// 			RadiusKm: 10,
// 		},
// 	}
// 	resp, err := s.service.GetNearbyFires(s.ctx, req)

// 	// Verify
// 	s.NoError(err)
// 	s.NotNil(resp)
// 	s.Len(resp.Msg.Wildfires, 1)
// 	s.Equal("fire-nearby", resp.Msg.Wildfires[0].FireId)
// 	s.mockRepo.AssertExpectations(s.T())
// }

// func (s *WildfireServiceTestSuite) TestProcessScrape_Success() {
// 	now := time.Now()
//     detections := []*scrapers.WildfireDetection{
//         {
//             ExternalID:    "fire-001",
//             Source:       "NASA",
//             Latitude:     34.0,
//             Longitude:    -118.5,
//             DiscoveredAt: now,
//         },
//     }

//     weatherData := &scrapers.WeatherData{
//         Temperature:   25.0,
//         WindSpeed:     10.0,
//         WindDirection: 180.0,
//         Humidity:      40.0,
//     }

//     // Use mock.Anything for context to accept any context type
//     s.mockScraper.On("Scrape", mock.Anything).Return(detections, nil)
//     s.mockWeather.On("GetWeatherData", mock.Anything, 34.0, -118.5).Return(weatherData, nil)
//     s.mockRepo.On("BulkUpsertFromDetections", mock.Anything, mock.Anything).Return(nil)

//     s.service.processScrape(s.ctx, s.mockScraper)

//     s.mockScraper.AssertExpectations(s.T())
//     s.mockWeather.AssertExpectations(s.T())
//     s.mockRepo.AssertExpectations(s.T())
// }

// func (s *WildfireServiceTestSuite) TestUpsertWildfire_Success() {
// 	now := time.Now().Format(time.RFC3339)
//     req := &connect.Request[nrtv1.UpsertRequest]{
//         Msg: &nrtv1.UpsertRequest{
//             Wildfire: &nrtv1.WildfireNRT{
//                 FireId:         "fire-123",
//                 Source:         "NASA",
//                 ExternalId:     "ext-123",
//                 Status:         "active",
//                 DiscoveryTime:  now,
//                 LastUpdated:    now,
//                 Location: &nrtv1.GeoPoint{
//                     Type:        "Point",
//                     Coordinates: []float64{-118.5, 34.0},
//                 },
//                 Weather: &nrtv1.WeatherData{
//                     Temperature:   25.0,
//                     WindSpeed:     10.0,
//                     WindDirection: 180.0,
//                     Humidity:      40.0,
//                     LastUpdated:   now,
//                 },
//                 Metadata: map[string]any{"key": "value"},
//             },
//         },
//     }

//     s.mockRepo.On("UpsertFromDetection", mock.Anything, mock.Anything).Return(nil)

//     resp, err := s.service.UpsertFromDetection(s.ctx, req)

//     s.NoError(err)
//     s.NotNil(resp)
//     s.NotEmpty(resp.Msg.DocId)
//     s.mockRepo.AssertExpectations(s.T())
// }

// func TestWildfireServiceSuite(t *testing.T) {
// 	suite.Run(t, new(WildfireServiceTestSuite))
// }