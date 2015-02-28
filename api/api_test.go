package api

import (
	"fmt"
	"github.com/martinp/atbapi/atb"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMarshal(t *testing.T) {
	data := struct{ Foo string }{"bar"}
	b, err := marshal(data, true)
	if err != nil {
		t.Fatal(err)
	}
	actual := string(b)
	expected := "{\n  \"Foo\": \"bar\"\n}"
	if actual != expected {
		t.Fatalf("Expected '%s', got '%s'", expected, actual)
	}
	b, err = marshal(data, false)
	if err != nil {
		t.Fatal(err)
	}
	actual = string(b)
	expected = "{\"Foo\":\"bar\"}"
	if actual != expected {
		t.Fatalf("Expected '%s', got '%s'", expected, actual)
	}
}

func newTestServer(path string, body string) *httptest.Server {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/soap+xml; charset=utf-8")
		fmt.Fprint(w, body)
	}
	mux := http.NewServeMux()
	mux.HandleFunc(path, handler)
	return httptest.NewServer(mux)
}

func TestGetBusStops(t *testing.T) {
	server := newTestServer("/", busStopsResponse)
	defer server.Close()
	atb := atb.Client{URL: server.URL}
	api := New(atb)
	_, err := api.getBusStops()
	if err != nil {
		t.Fatal(err)
	}
	cached, ok := api.cache.Get("stops")
	if !ok {
		t.Fatal("Expected true")
	}
	busStops, ok := cached.(BusStops)
	if !ok {
		t.Fatal("Expected true")
	}
	if len(busStops.Stops) != 1 {
		t.Fatal("Expected length to be 1")
	}
	if len(busStops.nodeIds) != 1 {
		t.Fatal("Expected length to be 1")
	}
}

func TestGetDepartures(t *testing.T) {
	server := newTestServer("/", forecastResponse)
	defer server.Close()
	atb := atb.Client{URL: server.URL}
	api := New(atb)
	_, err := api.getDepartures(16011376)
	if err != nil {
		t.Fatal(err)
	}
	cached, ok := api.cache.Get("16011376")
	if !ok {
		t.Fatal("Expected true")
	}
	departures, ok := cached.(Departures)
	if !ok {
		t.Fatal("Expected true")
	}
	if len(departures.Departures) != 1 {
		t.Fatal("Expected length to be 1")
	}
}

const busStopsResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap12:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap12="http://www.w3.org/2003/05/soap-envelope">
  <soap12:Body>
    <GetBusStopsListResponse xmlns="http://miz.it/infotransit">
      <GetBusStopsListResult>
{
  "Fermate": [
    {
      "cinAzienda": 1,
      "nomeAzienda": "AtB",
      "cinFermata": 100633,
      "codAzNodo": "16011376",
      "descrizione": "Prof. Brochs gt",
      "lon": "1157514",
      "lat": 9202874,
      "name": "Prof.",
      "codeMobile": "16011376 (Prof.)",
      "nomeMobile": "Prof. (16011376)"
    }
  ]
}
      </GetBusStopsListResult>
    </GetBusStopsListResponse>
  </soap12:Body>
</soap12:Envelope>`

const forecastResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap12:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap12="http://www.w3.org/2003/05/soap-envelope">
  <soap12:Body>
    <getUserRealTimeForecastByStopResponse xmlns="http://miz.it/infotransit">
      <getUserRealTimeForecastByStopResult>
{
  "total": 1,
  "timeServer": "2015-02-26 18:37",
  "InfoNodo": [
    {
      "nome_Az": "AtB",
      "codAzNodo": "16011376",
      "nomeNodo": "Prof.",
      "descrNodo": "Prof. Brochs gt",
      "bitMaskProprieta": "0",
      "codeMobile": "Prof. Brochs gt",
      "coordLon": "10.398126",
      "coordLat": "63.415535"
    }
  ],
  "Orari": [
    {
      "codAzLinea": "6",
      "descrizioneLinea": "6",
      "orario": "26.02.2015 18:38",
      "orarioSched": "26.02.2015 18:01",
      "statoPrevisione": "Prev",
      "capDest": "Munkegata M5",
      "turnoMacchina": "57",
      "descrizionePercorso": "39"
    }
  ]
}
      </getUserRealTimeForecastByStopResult>
    </getUserRealTimeForecastByStopResponse>
  </soap12:Body>
</soap12:Envelope>`
