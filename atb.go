package main

import (
	"bytes"
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

func (a *Atb) GetBusStops() ([]byte, error) {
	data := struct {
		Username string
		Password string
	}{a.Username, a.Password}
	method := a.Methods.GetBusStopsList
	json, err := a.post(method, data)
	if err != nil {
		return nil, err
	}
	return json, nil
}
