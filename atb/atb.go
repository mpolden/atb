package atb

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const URL = "http://st.atb.no/InfoTransit/userservices.asmx"

type Client struct {
	Username string
	Password string
	methods  Methods
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
	client.methods = newMethods()
	return client, nil
}

func (c *Client) post(m Method, data interface{}) ([]byte, error) {
	req, err := m.NewRequest(data)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBufferString(req)
	resp, err := http.Post(URL, "application/soap+xml", buf)
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
	method := c.methods.GetBusStopsList

	jsonBlob, err := c.post(method, values)
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
	method := c.methods.GetRealTimeForecast

	jsonBlob, err := c.post(method, values)
	if err != nil {
		return Forecasts{}, err
	}

	var forecasts Forecasts
	if err := json.Unmarshal(jsonBlob, &forecasts); err != nil {
		return Forecasts{}, err
	}
	return forecasts, nil
}
