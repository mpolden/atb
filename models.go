package main

import (
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

type busStops struct {
	Stops []busStop `json:"Fermate"`
}

type busStop struct {
	CompanyId   int    `json:"cinAzienda"`
	CompanyName string `json:"nomeAzienda"`
	StopId      int    `json:"cinFermata"`
	NodeId      string `json:"codAzNodo"`
	Description string `json:"descrizione"`
	Longitude   string `json:"lon"`
	Latitude    int    `json:"lat"`
	MobileCode  string `json:"codeMobile"`
	MobileName  string `json:"nomeMobile"`
}

func (s *busStops) Convert() (BusStops, error) {
	stops := make([]BusStop, 0, len(s.Stops))
	for _, stop := range s.Stops {
		converted, err := stop.Convert()
		if err != nil {
			return BusStops{}, err
		}
		stops = append(stops, converted)
	}
	return BusStops{Stops: stops}, nil
}

func (s *busStop) Convert() (BusStop, error) {
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
