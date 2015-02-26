package atb

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func newTestServer(path string, body string) *httptest.Server {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, body)
	}
	mux := http.NewServeMux()
	mux.HandleFunc(path, handler)
	return httptest.NewServer(mux)
}

func TestGetBusStops(t *testing.T) {
	server := newTestServer("/", busStopsResponse)
	defer server.Close()
	atb := Client{URL: server.URL}
	expected := BusStops{
		Stops: []BusStop{
			BusStop{
				StopId:      100633,
				NodeId:      "16011376",
				Description: "Prof. Brochs gt",
				Longitude:   "1157514",
				Latitude:    9202874,
				MobileCode:  "16011376 (Prof.)",
				MobileName:  "Prof. (16011376)",
			},
		},
	}
	stops, err := atb.GetBusStops()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(stops, expected) {
		t.Fatalf("Expected %+v, got %+v", expected, stops)
	}
}

func TestGetRealTimeForecast(t *testing.T) {
	server := newTestServer("/", forecastResponse)
	defer server.Close()
	atb := Client{URL: server.URL}
	forecasts, err := atb.GetRealTimeForecast(16011376)
	expected := Forecasts{
		Total: 1,
		Nodes: []NodeInfo{
			NodeInfo{
				Name:              "AtB",
				NodeId:            "16011376",
				NodeName:          "Prof.",
				NodeDescription:   "Prof. Brochs gt",
				BitMaskProperties: "0",
				Longitude:         "10.398126",
				Latitude:          "63.415535",
				MobileCode:        "Prof. Brochs gt",
			},
		},
		Forecasts: []Forecast{
			Forecast{
				LineId:                  "6",
				LineDescription:         "6",
				RegisteredDepartureTime: "26.02.2015 18:38",
				ScheduledDepartureTime:  "26.02.2015 18:01",
				StationForecast:         "Prev",
				Destination:             "Munkegata M5",
			},
		},
	}
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(forecasts, expected) {
		t.Fatalf("Expected %+v, got %+v", expected, forecasts)
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
