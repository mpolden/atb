package http

import (
	"github.com/mpolden/atb/entur"
)

// BusStops represents a list of bus stops.
type BusStops struct {
	Stops   []BusStop `json:"stops"`
	nodeIDs map[int]*BusStop
}

// BusStop represents a single bus stop.
type BusStop struct {
	URL         string  `json:"url"`
	StopID      int     `json:"stopId"`
	NodeID      int     `json:"nodeId"`
	Description string  `json:"description"`
	Longitude   float64 `json:"longitude"`
	Latitude    float64 `json:"latitude"`
	MobileCode  string  `json:"mobileCode"`
	MobileName  string  `json:"mobileName"`
}

// Departures represents a list of departures, from a given bus stop.
type Departures struct {
	URL            string      `json:"url"`
	TowardsCentrum *bool       `json:"isGoingTowardsCentrum,omitempty"`
	Departures     []Departure `json:"departures"`
}

// Departure represents a single departure in a given direction.
type Departure struct {
	LineID                  string `json:"line"`
	RegisteredDepartureTime string `json:"registeredDepartureTime,omitempty"`
	ScheduledDepartureTime  string `json:"scheduledDepartureTime"`
	Destination             string `json:"destination"`
	IsRealtimeData          bool   `json:"isRealtimeData"`
	TowardsCentrum          *bool  `json:"isGoingTowardsCentrum,omitempty"`
}

// Error represents an error in the API, which is returned to the user.
type Error struct {
	err     error
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func convertDepartures(enturDepartures []entur.Departure) Departures {
	departures := make([]Departure, 0, len(enturDepartures))
	const timeLayout = "2006-01-02T15:04:05.000"
	for _, d := range enturDepartures {
		scheduledDepartureTime := d.ScheduledDepartureTime.Format(timeLayout)
		registeredDepartureTime := ""
		if !d.RegisteredDepartureTime.IsZero() {
			registeredDepartureTime = d.RegisteredDepartureTime.Format(timeLayout)
		}
		towardsCentrum := d.Inbound
		departure := Departure{
			LineID:                  d.Line,
			ScheduledDepartureTime:  scheduledDepartureTime,
			RegisteredDepartureTime: registeredDepartureTime,
			Destination:             d.Destination,
			IsRealtimeData:          d.IsRealtime,
			TowardsCentrum:          &towardsCentrum,
		}
		departures = append(departures, departure)
	}
	return Departures{
		Departures: departures,
	}
}
