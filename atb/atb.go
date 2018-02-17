package atb

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
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

func (c *Client) post(r request, data interface{}) ([]byte, error) {
	body, err := r.Body(data)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(c.URL, "application/soap+xml", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := r.Decode(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// BusStops retrieves bus stops from AtBs API.
func (c *Client) BusStops() (BusStops, error) {
	values := struct {
		Username string
		Password string
	}{c.Username, c.Password}

	res, err := c.post(busStops, values)
	if err != nil {
		return BusStops{}, err
	}

	var stops BusStops
	if err := json.Unmarshal(res, &stops); err != nil {
		return BusStops{}, err
	}
	return stops, nil
}

// Forecasts retrieves forecasts from AtBs API, using nodeID to identify the bus stop.
func (c *Client) Forecasts(nodeID int) (Forecasts, error) {
	values := struct {
		Username string
		Password string
		NodeID   int
	}{c.Username, c.Password, nodeID}

	res, err := c.post(forecast, values)
	if err != nil {
		return Forecasts{}, err
	}

	var forecasts Forecasts
	if err := json.Unmarshal(res, &forecasts); err != nil {
		return Forecasts{}, err
	}
	return forecasts, nil
}
