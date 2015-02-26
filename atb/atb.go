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

	var stops busStops
	if err := json.Unmarshal(jsonBlob, &stops); err != nil {
		return BusStops{}, err
	}
	return stops.Convert()
}
