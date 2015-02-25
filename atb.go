package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const URL = "http://st.atb.no/InfoTransit/userservices.asmx"

type Atb struct {
	Username string
	Password string
	Methods  Methods
}

func (a *Atb) post(m Method, data interface{}) ([]byte, error) {
	req, err := m.CompileRequest(data)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBufferString(req)
	resp, err := http.Post(URL, "application/soap+xml", buf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	jsonResponse, err := m.ParseResponse(resp.Body)
	if err != nil {
		return nil, err
	}
	return jsonResponse, nil
}

func (a *Atb) GetBusStops() (BusStops, error) {
	values := struct {
		Username string
		Password string
	}{a.Username, a.Password}
	method := a.Methods.GetBusStopsList

	jsonBlob, err := a.post(method, values)
	if err != nil {
		return BusStops{}, err
	}

	var stops busStops
	if err := json.Unmarshal(jsonBlob, &stops); err != nil {
		return BusStops{}, err
	}
	return stops.Convert()
}
