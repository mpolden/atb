package atb

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// DefaultURL is the default AtB API URL.
const DefaultURL = "http://st.atb.no/New/InfoTransit/UserServices.asmx"

// Client represents a client which communicates with AtBs API.
type Client struct {
	Username string
	Password string
	URL      string
}

// BusStops represents a list of bus stops.
type BusStops struct {
	Stops []BusStop `json:"Fermate"`
}

// BusStop represents a bus stop.
type BusStop struct {
	StopID      int    `json:"cinFermata"`
	NodeID      string `json:"codAzNodo"`
	Description string `json:"descrizione"`
	Longitude   string `json:"lon"`
	Latitude    int    `json:"lat"`
	MobileCode  string `json:"codeMobile"`
	MobileName  string `json:"nomeMobile"`
}

// Forecasts represents a list of forecasts.
type Forecasts struct {
	Nodes     []NodeInfo `json:"InfoNodo"`
	Forecasts []Forecast `json:"Orari"`
	Total     int        `json:"total"`
}

// NodeInfo represents a bus stop, returned as a part of a forecast.
type NodeInfo struct {
	Name              string `json:"nome_Az"`
	NodeID            string `json:"codAzNodo"`
	NodeName          string `json:"nomeNodo"`
	NodeDescription   string `json:"descrNodo"`
	BitMaskProperties string `json:"bitMaskProprieta"`
	MobileCode        string `json:"codeMobile"`
	Longitude         string `json:"coordLon"`
	Latitude          string `json:"coordLat"`
}

// Forecast represents a single forecast.
type Forecast struct {
	LineID                  string `json:"codAzLinea"`
	LineDescription         string `json:"descrizioneLinea"`
	RegisteredDepartureTime string `json:"orario"`
	ScheduledDepartureTime  string `json:"orarioSched"`
	StationForecast         string `json:"statoPrevisione"`
	Destination             string `json:"capDest"`
}

type busStopsRequest struct {
	XMLName xml.Name `xml:"Envelope"`
	Result  []byte   `xml:"Body>GetBusStopsListResponse>GetBusStopsListResult"`
}

type forecastRequest struct {
	XMLName xml.Name `xml:"Envelope"`
	Result  []byte   `xml:"Body>getUserRealTimeForecastByStopResponse>getUserRealTimeForecastByStopResult"`
}

// NewFromConfig creates a new client where name is the path to the config file.
func NewFromConfig(name string) (Client, error) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return Client{}, err
	}
	var client Client
	if err := json.Unmarshal(data, &client); err != nil {
		return Client{}, err
	}
	if client.URL == "" {
		client.URL = DefaultURL
	}
	return client, nil
}

func (c *Client) postXML(body string, dst interface{}) error {
	resp, err := http.Post(c.URL, "application/soap+xml", strings.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	dec := xml.NewDecoder(resp.Body)
	if err := dec.Decode(&dst); err != nil {
		return err
	}
	return nil
}

// BusStops retrieves bus stops from AtBs API.
func (c *Client) BusStops() (BusStops, error) {
	req := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<soap12:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap12="http://www.w3.org/2003/05/soap-envelope">
  <soap12:Body>
    <GetBusStopsList xmlns="http://miz.it/infotransit">
      <auth>
        <user>%s</user>
        <password>%s</password>
      </auth>
    </GetBusStopsList>
  </soap12:Body>
</soap12:Envelope>`, c.Username, c.Password)

	var stopsRequest busStopsRequest
	if err := c.postXML(req, &stopsRequest); err != nil {
		return BusStops{}, err
	}

	var stops BusStops
	if err := json.Unmarshal(stopsRequest.Result, &stops); err != nil {
		return BusStops{}, err
	}
	return stops, nil
}

// Forecasts retrieves forecasts from AtBs API, using nodeID to identify the bus stop.
func (c *Client) Forecasts(nodeID int) (Forecasts, error) {
	req := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<soap12:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap12="http://www.w3.org/2003/05/soap-envelope">
  <soap12:Body>
    <getUserRealTimeForecastByStop xmlns="http://miz.it/infotransit">
      <auth>
        <user>%s</user>
        <password>%s</password>
      </auth>
      <busStopId>%d</busStopId>
    </getUserRealTimeForecastByStop>
  </soap12:Body>
</soap12:Envelope>`, c.Username, c.Password, nodeID)

	var forecastRequest forecastRequest
	if err := c.postXML(req, &forecastRequest); err != nil {
		return Forecasts{}, err
	}

	var forecasts Forecasts
	if err := json.Unmarshal(forecastRequest.Result, &forecasts); err != nil {
		return Forecasts{}, err
	}
	return forecasts, nil
}
