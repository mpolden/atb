package atb

import (
	"bytes"
	"encoding/xml"
	"io"
	"text/template"
)

type request interface {
	Body(data interface{}) (io.Reader, error)
	Decode(r io.Reader) ([]byte, error)
}

func compileTemplate(t *template.Template, data interface{}) (io.Reader, error) {
	var b bytes.Buffer
	if err := t.Execute(&b, data); err != nil {
		return nil, err
	}
	return &b, nil
}

type busStopsRequest struct {
	XMLName  xml.Name `xml:"Envelope"`
	Result   []byte   `xml:"Body>GetBusStopsListResponse>GetBusStopsListResult"`
	template *template.Template
}

type forecastRequest struct {
	XMLName  xml.Name `xml:"Envelope"`
	Result   []byte   `xml:"Body>getUserRealTimeForecastByStopResponse>getUserRealTimeForecastByStopResult"`
	template *template.Template
}

func (m *busStopsRequest) Body(data interface{}) (io.Reader, error) {
	return compileTemplate(m.template, data)
}

func (m *busStopsRequest) Decode(r io.Reader) ([]byte, error) {
	var stops busStopsRequest
	dec := xml.NewDecoder(r)
	if err := dec.Decode(&stops); err != nil {
		return nil, err
	}
	return stops.Result, nil
}

func (m *forecastRequest) Body(data interface{}) (io.Reader, error) {
	return compileTemplate(m.template, data)
}

func (m *forecastRequest) Decode(r io.Reader) ([]byte, error) {
	var forecast forecastRequest
	dec := xml.NewDecoder(r)
	if err := dec.Decode(&forecast); err != nil {
		return nil, err
	}
	return forecast.Result, nil
}

func templateMust(src string) *template.Template {
	return template.Must(template.New("xml").Parse(src))
}

var (
	busStops = &busStopsRequest{
		template: templateMust(`<?xml version="1.0" encoding="utf-8"?>
<soap12:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap12="http://www.w3.org/2003/05/soap-envelope">
  <soap12:Body>
    <GetBusStopsList xmlns="http://miz.it/infotransit">
      <auth>
        <user>{{.Username}}</user>
        <password>{{.Password}}</password>
      </auth>
    </GetBusStopsList>
  </soap12:Body>
</soap12:Envelope>`),
	}
	forecast = &forecastRequest{
		template: templateMust(`<?xml version="1.0" encoding="utf-8"?>
<soap12:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap12="http://www.w3.org/2003/05/soap-envelope">
  <soap12:Body>
    <getUserRealTimeForecastByStop xmlns="http://miz.it/infotransit">
      <auth>
        <user>{{.Username}}</user>
        <password>{{.Password}}</password>
      </auth>
      <busStopId>{{.NodeID}}</busStopId>
    </getUserRealTimeForecastByStop>
  </soap12:Body>
</soap12:Envelope>`),
	}
)
