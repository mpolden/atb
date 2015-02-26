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
	GetBusStopsList Method
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

func newMethods() Methods {
	return Methods{
		GetBusStopsList: &GetBusStopsList{
			template: getBusStopsTemplate,
		},
	}
}
