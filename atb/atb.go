package atb

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const DefaultURL = "http://st.atb.no/InfoTransit/userservices.asmx"

type Client struct {
	Username string
	Password string
	URL      string
}

type BusStops struct {
	Stops []BusStop `json:"Fermate"`
}

type BusStop struct {
	StopId      int    `json:"cinFermata"`
	NodeId      string `json:"codAzNodo"`
	Description string `json:"descrizione"`
	Longitude   string `json:"lon"`
	Latitude    int    `json:"lat"`
	MobileCode  string `json:"codeMobile"`
	MobileName  string `json:"nomeMobile"`
}

type Forecasts struct {
	Nodes     []NodeInfo `json:"InfoNodo"`
	Forecasts []Forecast `json:"Orari"`
	Total     int        `json:"total"`
}

type NodeInfo struct {
	Name              string `json:"nome_Az"`
	NodeId            string `json:"codAzNodo"`
	NodeName          string `json:"nomeNodo"`
	NodeDescription   string `json:"descrNodo"`
	BitMaskProperties string `json:"bitMaskProprieta"`
	MobileCode        string `json:"codeMobile"`
	Longitude         string `json:"coordLon"`
	Latitude          string `json:"coordLat"`
}

type Forecast struct {
	LineId                  string `json:"codAzLinea"`
	LineDescription         string `json:"descrizioneLinea"`
	RegisteredDepartureTime string `json:"orario"`
	ScheduledDepartureTime  string `json:"orarioSched"`
	StationForecast         string `json:"statoPrevisione"`
	Destination             string `json:"capDest"`
}

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

func (c *Client) post(m Method, data interface{}) ([]byte, error) {
	req, err := m.NewRequest(data)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBufferString(req)
	resp, err := http.Post(c.URL, "application/soap+xml", buf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	jsonBlob, err := m.ParseResponse(body)
	if err != nil {
		return nil, err
	}
	return jsonBlob, nil
}

func (c *Client) GetBusStops() (BusStops, error) {
	values := struct {
		Username string
		Password string
	}{c.Username, c.Password}

	jsonBlob, err := c.post(getBusStopsList, values)
	if err != nil {
		return BusStops{}, err
	}

	var stops BusStops
	if err := json.Unmarshal(jsonBlob, &stops); err != nil {
		return BusStops{}, err
	}
	return stops, nil
}

func (c *Client) GetRealTimeForecast(nodeId int) (Forecasts, error) {
	values := struct {
		Username string
		Password string
		NodeId   int
	}{c.Username, c.Password, nodeId}

	jsonBlob, err := c.post(getRealTimeForecast, values)
	if err != nil {
		return Forecasts{}, err
	}

	var forecasts Forecasts
	if err := json.Unmarshal(jsonBlob, &forecasts); err != nil {
		return Forecasts{}, err
	}
	return forecasts, nil
}
