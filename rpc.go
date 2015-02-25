package main

import (
	"bytes"
	"fmt"
	"gopkg.in/xmlpath.v1"
	"io"
	"io/ioutil"
	"path/filepath"
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

func Parse(tmplFile string) (*template.Template, error) {
	contents, err := ioutil.ReadFile(tmplFile)
	if err != nil {
		return nil, err
	}
	tmpl, err := template.New("").Parse(string(contents))
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

func NewMethods(path string) (Methods, error) {
	getBusStopsPath := xmlpath.MustCompile("/Envelope/Body/" +
		"GetBusStopsListResponse/GetBusStopsListResult")
	getBusStopsTmpl, err := Parse(filepath.Join(path,
		"GetBusStopsList.tmpl"))
	if err != nil {
		return Methods{}, err
	}
	return Methods{
		GetBusStopsList: Method{
			Template: getBusStopsTmpl,
			Path:     getBusStopsPath,
		},
	}, nil
}
