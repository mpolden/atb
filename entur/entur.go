package entur

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// DefaultURL is the default Entur Journey Planner API URL. Documentation at
// https://developer.entur.org/pages-journeyplanner-journeyplanner-v2.
const DefaultURL = "https://api.entur.io/journey-planner/v2/graphql"

// Client implements a client for the Entur Journey Planner API.
type Client struct{ URL string }

// New creates a new client using the API found at url.
func New(url string) *Client {
	if url == "" {
		url = DefaultURL
	}
	return &Client{URL: url}
}

// Departure represents a bus departure from a stop.
type Departure struct {
	Line                    string
	RegisteredDepartureTime time.Time
	ScheduledDepartureTime  time.Time
	Destination             string
	IsRealtime              bool
	Inbound                 bool
}

type response struct {
	Data data `json:"data"`
}

type data struct {
	StopPlace stopPlace `json:"stopPlace"`
}

type stopPlace struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	EstimatedCalls []estimatedCall `json:"estimatedCalls"`
}

type estimatedCall struct {
	Realtime              bool               `json:"realtime"`
	ExpectedDepartureTime string             `json:"expectedDepartureTime"`
	ActualDepartureTime   string             `json:"actualDepartureTime"`
	DestinationDisplay    destinationDisplay `json:"destinationDisplay"`
	ServiceJourney        serviceJourney     `json:"serviceJourney"`
}

type destinationDisplay struct {
	FrontText string `json:"frontText"`
}

type serviceJourney struct {
	JourneyPattern journeyPattern `json:"journeyPattern"`
}

type journeyPattern struct {
	DirectionType string `json:"directionType"`
	Line          line   `json:"line"`
}

type line struct {
	PublicCode string `json:"publicCode"`
}

// Departures returns departures from the given stop ID. Use https://stoppested.entur.org/ to determine stop IDs.
func (c *Client) Departures(count, stopID int) ([]Departure, error) {
	// https://api.entur.io/journey-planner/v2/ide/ for query testing
	query := fmt.Sprintf(`{"query":"{stopPlace(id:\"NSR:StopPlace:%d\"){id name estimatedCalls(numberOfDepartures:%d){realtime expectedDepartureTime actualDepartureTime destinationDisplay{frontText}serviceJourney{journeyPattern{directionType line{publicCode}}}}}}"}`, stopID, count)
	req, err := http.NewRequest("POST", c.URL, strings.NewReader(query))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	// Identify this client. See https://developer.entur.org/pages-journeyplanner-journeyplanner-v2
	req.Header.Set("ET-Client-Name", "github_mpolden-atb")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	json, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return parseDepartures(json)
}

func parseDepartures(jsonData []byte) ([]Departure, error) {
	var r response
	if err := json.Unmarshal(jsonData, &r); err != nil {
		return nil, err
	}
	const timeLayout = "2006-01-02T15:04:05-0700"
	departures := make([]Departure, 0, len(r.Data.StopPlace.EstimatedCalls))
	for _, ec := range r.Data.StopPlace.EstimatedCalls {
		scheduledDepartureTime, err := time.Parse(timeLayout, ec.ExpectedDepartureTime)
		if err != nil {
			return nil, err
		}
		registeredDepartureTime := time.Time{}
		if ec.ActualDepartureTime != "" {
			t, err := time.Parse(timeLayout, ec.ActualDepartureTime)
			if err != nil {
				return nil, err
			}
			registeredDepartureTime = t
		}
		inbound := ec.ServiceJourney.JourneyPattern.DirectionType == "inbound"
		d := Departure{
			Line:                    ec.ServiceJourney.JourneyPattern.Line.PublicCode,
			RegisteredDepartureTime: registeredDepartureTime,
			ScheduledDepartureTime:  scheduledDepartureTime,
			Destination:             ec.DestinationDisplay.FrontText,
			IsRealtime:              ec.Realtime,
			Inbound:                 inbound,
		}
		departures = append(departures, d)
	}
	return departures, nil
}
