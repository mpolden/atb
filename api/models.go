package api

import (
	"github.com/martinp/atbapi/atb"
	"strconv"
)

type BusStops struct {
	Stops []BusStop `json:"stops"`
}

type BusStop struct {
	StopId      int     `json:"stopId"`
	NodeId      int     `json:"nodeId"`
	Description string  `json:"description"`
	Longitude   float64 `json:"longitude"`
	Latitude    float64 `json:"latitude"`
	MobileCode  string  `json:"mobileCode"`
	MobileName  string  `json:"mobileName"`
}

func convertBusStop(s atb.BusStop) (BusStop, error) {
	nodeId, err := strconv.Atoi(s.NodeId)
	if err != nil {
		return BusStop{}, err
	}
	longitude, err := strconv.ParseFloat(s.Longitude, 64)
	if err != nil {
		return BusStop{}, err
	}
	latitude := float64(s.Latitude)
	return BusStop{
		StopId:      s.StopId,
		NodeId:      nodeId,
		Description: s.Description,
		Longitude:   longitude,
		Latitude:    latitude,
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
