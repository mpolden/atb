package api

import (
	"github.com/martinp/atbapi/atb"
	"strconv"
	"strings"
	"time"
)

// BusStops represents a list of bus stops.
type BusStops struct {
	Stops   []BusStop `json:"stops"`
	nodeIDs map[int]struct{}
}

// BusStop represents a single bus stop.
type BusStop struct {
	StopID      int    `json:"stopId"`
	NodeID      int    `json:"nodeId"`
	Description string `json:"description"`
	Longitude   int    `json:"longitude"`
	Latitude    int    `json:"latitude"`
	MobileCode  string `json:"mobileCode"`
	MobileName  string `json:"mobileName"`
}

// Departures represents a list of departures, from a given bus stop.
type Departures struct {
	TowardsCentrum bool        `json:"isGoingTowardsCentrum"`
	Departures     []Departure `json:"departures"`
}

// Departure represents a single departure in a given direction.
type Departure struct {
	LineID                  string `json:"line"`
	RegisteredDepartureTime string `json:"registeredDepartureTime"`
	ScheduledDepartureTime  string `json:"scheduledDepartureTime"`
	Destination             string `json:"destination"`
	IsRealtimeData          bool   `json:"isRealtimeData"`
}

// Error represents an error in the API, which is returned to the user.
type Error struct {
	err     error
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func convertBusStop(s atb.BusStop) (BusStop, error) {
	nodeID, err := strconv.Atoi(s.NodeID)
	if err != nil {
		return BusStop{}, err
	}
	longitude, err := strconv.Atoi(s.Longitude)
	if err != nil {
		return BusStop{}, err
	}
	return BusStop{
		StopID:      s.StopID,
		NodeID:      nodeID,
		Description: s.Description,
		Longitude:   longitude,
		Latitude:    s.Latitude,
		MobileCode:  s.MobileCode,
		MobileName:  s.MobileName,
	}, nil
}

func convertBusStops(s atb.BusStops) (BusStops, error) {
	stops := make([]BusStop, 0, len(s.Stops))
	for _, stop := range s.Stops {
		converted, err := convertBusStop(stop)
		if err != nil {
			return BusStops{}, err
		}
		stops = append(stops, converted)
	}
	return BusStops{Stops: stops}, nil
}

// ConvertTime converts time from AtBs format to ISO 8601.
func ConvertTime(src string) (string, error) {
	t, err := time.Parse("02.01.2006 15:04", src)
	if err != nil {
		return "", err
	}
	return t.Format("2006-01-02T15:04:05.000"), nil
}

// IsRealtime returns a boolean indicating whether stationForecast is realtime.
func IsRealtime(stationForecast string) bool {
	return strings.EqualFold(stationForecast, "prev")
}

func convertForecast(f atb.Forecast) (Departure, error) {
	registeredDeparture, err := ConvertTime(f.RegisteredDepartureTime)
	if err != nil {
		return Departure{}, err
	}
	scheduledDeparture, err := ConvertTime(f.ScheduledDepartureTime)
	if err != nil {
		return Departure{}, err
	}
	return Departure{
		LineID:                  f.LineID,
		Destination:             f.Destination,
		RegisteredDepartureTime: registeredDeparture,
		ScheduledDepartureTime:  scheduledDeparture,
		IsRealtimeData:          IsRealtime(f.StationForecast),
	}, nil
}

// IsTowardsCentrum returns a boolean indicating whether a bus stop, identified
// by nodeID, is going to the centrum
func IsTowardsCentrum(nodeID int) bool {
	return (nodeID/1000)%2 == 1
}

func convertForecasts(f atb.Forecasts) (Departures, error) {
	towardsCentrum := false
	if len(f.Nodes) > 0 {
		nodeID, err := strconv.Atoi(f.Nodes[0].NodeID)
		if err != nil {
			return Departures{}, err
		}
		towardsCentrum = IsTowardsCentrum(nodeID)
	}
	departures := make([]Departure, 0, len(f.Forecasts))
	for _, forecast := range f.Forecasts {
		departure, err := convertForecast(forecast)
		if err != nil {
			return Departures{}, err
		}
		departures = append(departures, departure)
	}
	return Departures{
		TowardsCentrum: towardsCentrum,
		Departures:     departures,
	}, nil
}
