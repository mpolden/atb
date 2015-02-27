package api

import (
	"github.com/martinp/atbapi/atb"
	"strconv"
	"strings"
	"time"
)

type BusStops struct {
	Stops   []BusStop `json:"stops"`
	nodeIds map[int]struct{}
}

type BusStop struct {
	StopId      int    `json:"stopId"`
	NodeId      int    `json:"nodeId"`
	Description string `json:"description"`
	Longitude   int    `json:"longitude"`
	Latitude    int    `json:"latitude"`
	MobileCode  string `json:"mobileCode"`
	MobileName  string `json:"mobileName"`
}

type Departures struct {
	TowardsCentrum bool        `json:"isGoingTowardsCentrum"`
	Departures     []Departure `json:"departures"`
}

type Departure struct {
	LineId                  string `json:"line"`
	RegisteredDepartureTime string `json:"registeredDepartureTime"`
	ScheduledDepartureTime  string `json:"scheduledDepartureTime"`
	Destination             string `json:"destination"`
	IsRealtimeData          bool   `json:"isRealtimeData"`
}

type Error struct {
	error   error  `json:"-"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func convertBusStop(s atb.BusStop) (BusStop, error) {
	nodeId, err := strconv.Atoi(s.NodeId)
	if err != nil {
		return BusStop{}, err
	}
	longitude, err := strconv.Atoi(s.Longitude)
	if err != nil {
		return BusStop{}, err
	}
	return BusStop{
		StopId:      s.StopId,
		NodeId:      nodeId,
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

func ConvertTime(src string) (string, error) {
	t, err := time.Parse("02.01.2006 15:04", src)
	if err != nil {
		return "", err
	}
	return t.Format("2006-01-02T15:04:05.000"), nil
}

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
		LineId:                  f.LineId,
		Destination:             f.Destination,
		RegisteredDepartureTime: registeredDeparture,
		ScheduledDepartureTime:  scheduledDeparture,
		IsRealtimeData:          IsRealtime(f.StationForecast),
	}, nil
}

func IsTowardsCentrum(nodeId int) bool {
	return (nodeId/1000)%2 == 1
}

func convertForecasts(f atb.Forecasts) (Departures, error) {
	towardsCentrum := false
	if len(f.Nodes) > 0 {
		nodeId, err := strconv.Atoi(f.Nodes[0].NodeId)
		if err != nil {
			return Departures{}, err
		}
		towardsCentrum = IsTowardsCentrum(nodeId)
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
