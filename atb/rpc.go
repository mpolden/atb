package atb

import (
	"bytes"
	"fmt"
	"gopkg.in/xmlpath.v1"
	"io"
	"text/template"
)

type Method struct {
	Template *template.Template
	Path     *xmlpath.Path
}

type Methods struct {
	GetBusStopsList Method
}

func (m *Method) CompileRequest(data interface{}) (string, error) {
	var b bytes.Buffer
	if err := m.Template.Execute(&b, data); err != nil {
		return "", err
	}
	return b.String(), nil
}

func (m *Method) ParseResponse(reader io.Reader) ([]byte, error) {
	node, err := xmlpath.Parse(reader)
	if err != nil {
		return nil, err
	}
	value, ok := m.Path.Bytes(node)
	if !ok {
		return nil, fmt.Errorf("could not find node")
	}
	return value, nil
}

func createMethods() Methods {
	getBusStopsPath := xmlpath.MustCompile("/Envelope/Body/" +
		"GetBusStopsListResponse/GetBusStopsListResult")
	return Methods{
		GetBusStopsList: Method{
			Template: GetBusStopsTemplate,
			Path:     getBusStopsPath,
		},
	}
}
