package atb

import (
	"bytes"
	"encoding/xml"
	"text/template"
)

type method interface {
	NewRequest(data interface{}) (string, error)
	ParseResponse(body []byte) ([]byte, error)
}

func compileTemplate(t *template.Template, data interface{}) (string, error) {
	var b bytes.Buffer
	if err := t.Execute(&b, data); err != nil {
		return "", err
	}
	return b.String(), nil
}

type busStopsListMethod struct {
	XMLName  xml.Name `xml:"Envelope"`
	Result   []byte   `xml:"Body>GetBusStopsListResponse>GetBusStopsListResult"`
	template *template.Template
}

type realTimeForecastMethod struct {
	XMLName  xml.Name `xml:"Envelope"`
	Result   []byte   `xml:"Body>getUserRealTimeForecastByStopResponse>getUserRealTimeForecastByStopResult"`
	template *template.Template
}

func (m *busStopsListMethod) NewRequest(data interface{}) (string, error) {
	return compileTemplate(m.template, data)
}

func (m *busStopsListMethod) ParseResponse(body []byte) ([]byte, error) {
	var stops busStopsListMethod
	if err := xml.Unmarshal(body, &stops); err != nil {
		return nil, err
	}
	return stops.Result, nil
}

func (m *realTimeForecastMethod) NewRequest(data interface{}) (string, error) {
	return compileTemplate(m.template, data)
}

func (m *realTimeForecastMethod) ParseResponse(body []byte) ([]byte, error) {
	var forecast realTimeForecastMethod
	if err := xml.Unmarshal(body, &forecast); err != nil {
		return nil, err
	}
	return forecast.Result, nil
}

func templateMust(src string) *template.Template {
	return template.Must(template.New("xml").Parse(src))
}

var (
	busStopsList = &busStopsListMethod{
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
	realTimeForecast = &realTimeForecastMethod{
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
