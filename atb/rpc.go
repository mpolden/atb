package atb

import (
	"bytes"
	"encoding/xml"
	"text/template"
)

type Method interface {
	NewRequest(data interface{}) (string, error)
	ParseResponse(body []byte) ([]byte, error)
}

type Methods struct {
	GetBusStopsList     Method
	GetRealTimeForecast Method
}

type GetBusStopsList struct {
	XMLName  xml.Name           `xml:"Envelope"`
	Result   []byte             `xml:"Body>GetBusStopsListResponse>GetBusStopsListResult"`
	template *template.Template `xml:"-"`
}

func compileTemplate(t *template.Template, data interface{}) (string, error) {
	var b bytes.Buffer
	if err := t.Execute(&b, data); err != nil {
		return "", err
	}
	return b.String(), nil
}

func (m *GetBusStopsList) NewRequest(data interface{}) (string, error) {
	return compileTemplate(m.template, data)
}

func (m *GetBusStopsList) ParseResponse(body []byte) ([]byte, error) {
	var stops GetBusStopsList
	if err := xml.Unmarshal(body, &stops); err != nil {
		return nil, err
	}
	return stops.Result, nil
}

type GetRealTimeForecast struct {
	XMLName  xml.Name           `xml:"Envelope"`
	Result   []byte             `xml:"Body>getUserRealTimeForecastByStopResponse>getUserRealTimeForecastByStopResult"`
	template *template.Template `xml:"-"`
}

func (m *GetRealTimeForecast) NewRequest(data interface{}) (string, error) {
	return compileTemplate(m.template, data)
}

func (m *GetRealTimeForecast) ParseResponse(body []byte) ([]byte, error) {
	var forecast GetRealTimeForecast
	if err := xml.Unmarshal(body, &forecast); err != nil {
		return nil, err
	}
	return forecast.Result, nil
}

func newMethods() Methods {
	return Methods{
		GetBusStopsList: &GetBusStopsList{
			template: getBusStopsTemplate,
		},
		GetRealTimeForecast: &GetRealTimeForecast{
			template: getRealTimeForecastTemplate,
		},
	}
}
