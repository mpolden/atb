package atb

import (
	"text/template"
)

func templateMust(src string) *template.Template {
	return template.Must(template.New("xml").Parse(src))
}

var (
	getBusStopsTemplate = templateMust(`<?xml version="1.0" encoding="utf-8"?>
<soap12:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap12="http://www.w3.org/2003/05/soap-envelope">
  <soap12:Body>
    <GetBusStopsList xmlns="http://miz.it/infotransit">
      <auth>
        <user>{{.Username}}</user>
        <password>{{.Password}}</password>
      </auth>
    </GetBusStopsList>
  </soap12:Body>
</soap12:Envelope>`)
)
