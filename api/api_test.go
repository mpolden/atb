package api

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mpolden/atbapi/atb"
)

func atbTestServer() *httptest.Server {
	handler := func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		xml := string(b)
		w.Header().Set("Content-Type", "application/soap+xml; charset=utf-8")
		if strings.Contains(xml, "GetBusStopsList") {
			fmt.Fprint(w, busStopsResponse)
		} else if strings.Contains(xml, "getUserRealTimeForecastByStop") {
			fmt.Fprint(w, forecastResponse)
		} else {
			panic("unknown request body: " + xml)
		}
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	return httptest.NewServer(mux)
}

func testServers() (*httptest.Server, *httptest.Server) {
	atbServer := atbTestServer()
	atb := atb.Client{URL: atbServer.URL}
	api := New(atb, 168*time.Hour, 1*time.Minute, false)
	return atbServer, httptest.NewServer(api.Handler())
}

func httpGet(url string) (string, string, int, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", "", 0, err
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", "", 0, err
	}
	return string(data), res.Header["Content-Type"][0], res.StatusCode, nil
}

func TestAPI(t *testing.T) {
	atbServer, server := testServers()
	defer atbServer.Close()
	defer server.Close()
	log.SetOutput(ioutil.Discard)

	var tests = []struct {
		url      string
		response string
		status   int
	}{
		// Unknown resources
		{"/not-found", `{"status":404,"message":"Resource not found"}`, 404},
		// List know URLs
		{"/", fmt.Sprintf(`{"urls":["%s/api/v1/busstops","%s/api/v1/departures"]}`, server.URL, server.URL), 200},
		// List all bus stops
		{"/api/v1/busstops", fmt.Sprintf(`{"stops":[{"url":"%s/api/v1/busstops/16011376","stopId":100633,"nodeId":16011376,"description":"Prof. Brochs gt","longitude":10.398126,"latitude":63.415535,"mobileCode":"16011376 (Prof.)","mobileName":"Prof. (16011376)"}]}`, server.URL), 200},
		// List all departures
		{"/api/v1/departures", fmt.Sprintf(`{"urls":["%s/api/v1/departures/16011376"]}`, server.URL), 200},
		// Show specific bus stop
		{"/api/v1/busstops/", `{"status":400,"message":"Invalid nodeID"}`, 400},
		{"/api/v1/busstops/foo", `{"status":400,"message":"Invalid nodeID"}`, 400},
		{"/api/v1/busstops/42", `{"status":404,"message":"Unknown bus stop"}`, 404},
		{"/api/v1/busstops/16011376", fmt.Sprintf(`{"url":"%s/api/v1/busstops/16011376","stopId":100633,"nodeId":16011376,"description":"Prof. Brochs gt","longitude":10.398126,"latitude":63.415535,"mobileCode":"16011376 (Prof.)","mobileName":"Prof. (16011376)"}`, server.URL), 200},
		// Show specific departure
		{"/api/v1/departures/", `{"status":400,"message":"Invalid nodeID"}`, 400},
		{"/api/v1/departures/foo", `{"status":400,"message":"Invalid nodeID"}`, 400},
		{"/api/v1/departures/42", `{"status":404,"message":"Unknown bus stop"}`, 404},
		{"/api/v1/departures/16011376", fmt.Sprintf(`{"url":"%s/api/v1/departures/16011376","isGoingTowardsCentrum":true,"departures":[{"line":"6","registeredDepartureTime":"2015-02-26T18:38:00.000","scheduledDepartureTime":"2015-02-26T18:01:00.000","destination":"Munkegata M5","isRealtimeData":true}]}`, server.URL), 200},
	}
	for _, tt := range tests {
		data, contentType, status, err := httpGet(server.URL + tt.url)
		if err != nil {
			t.Fatal(err)
		}
		if contentType != "application/json" {
			t.Errorf("want content-type application/json for %s, got %s", tt.url, contentType)
		}
		if got := status; status != tt.status {
			t.Errorf("want status %d for %s, got %d", tt.status, tt.url, got)
		}
		if got := string(data); got != tt.response {
			t.Errorf("want response %s for %s, got %s", tt.response, tt.url, got)
		}
	}
}

func TestURLPrefix(t *testing.T) {
	var tests = []struct {
		in  *http.Request
		out string
	}{
		{&http.Request{Host: "foo"}, "http://foo"},
		{&http.Request{Host: "", RemoteAddr: "127.0.0.1"}, "http://127.0.0.1"},
		{&http.Request{Host: "bar", TLS: &tls.ConnectionState{}}, "https://bar"},
		{&http.Request{Host: "baz", Header: map[string][]string{"X-Forwarded-Proto": []string{"https"}}}, "https://baz"},
		{&http.Request{Host: "qux", Header: map[string][]string{"X-Forwarded-Proto": []string{}}}, "http://qux"},
	}
	for _, tt := range tests {
		prefix := urlPrefix(tt.in)
		if prefix != tt.out {
			t.Errorf("want %s, got %s", tt.out, prefix)
		}
	}
}

func TestGetBusStops(t *testing.T) {
	server := atbTestServer()
	defer server.Close()
	atb := atb.Client{URL: server.URL}
	api := New(atb, 168*time.Hour, 1*time.Minute, false)
	_, _, err := api.getBusStops("")
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
	if len(busStops.nodeIDs) != 1 {
		t.Fatal("Expected length to be 1")
	}
}

func TestGetBusStopsCache(t *testing.T) {
	server := atbTestServer()
	defer server.Close()
	atb := atb.Client{URL: server.URL}
	api := New(atb, 168*time.Hour, 1*time.Minute, false)
	_, hit, err := api.getBusStops("")
	if err != nil {
		t.Fatal(err)
	}
	if hit {
		t.Error("Expected false")
	}
	_, hit, err = api.getBusStops("")
	if err != nil {
		t.Fatal(err)
	}
	if !hit {
		t.Error("Expected true")
	}
}

func TestGetDepartures(t *testing.T) {
	server := atbTestServer()
	defer server.Close()
	atb := atb.Client{URL: server.URL}
	api := New(atb, 168*time.Hour, 1*time.Minute, false)
	_, _, err := api.getDepartures("", 16011376)
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

func TestGetDeparturesCache(t *testing.T) {
	server := atbTestServer()
	defer server.Close()
	atb := atb.Client{URL: server.URL}
	api := New(atb, 168*time.Hour, 1*time.Minute, false)
	_, hit, err := api.getDepartures("", 16011376)
	if err != nil {
		t.Fatal(err)
	}
	if hit {
		t.Error("Expected false")
	}
	_, hit, err = api.getDepartures("", 16011376)
	if err != nil {
		t.Fatal(err)
	}
	if !hit {
		t.Error("Expected true")
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
