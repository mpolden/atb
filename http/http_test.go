package http

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mpolden/atb/entur"
)

func apiTestServer() *httptest.Server {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprint(w, enturResponse)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	return httptest.NewServer(mux)
}

func testServers() (*httptest.Server, *Server) {
	apiServer := apiTestServer()
	entur := &entur.Client{URL: apiServer.URL}
	server := New(entur, 168*time.Hour, 1*time.Minute, false)
	return apiServer, server
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
	apiServer, server := testServers()
	httpSrv := httptest.NewServer(server.Handler())
	defer apiServer.Close()
	defer httpSrv.Close()
	log.SetOutput(ioutil.Discard)

	var tests = []struct {
		url      string
		response string
		status   int
	}{
		// Unknown resources
		{"/not-found", `{"status":404,"message":"Resource not found"}`, 404},
		// List know URLs
		{"/", fmt.Sprintf(`{"urls":["%s/api/v2/departures"]}`, httpSrv.URL), 200},
		// Show specific departure (v2)
		{"/api/v2/departures", `{"status":400,"message":"Invalid stop ID. Use https://stoppested.entur.org/ to find stop IDs."}`, 400},
		{"/api/v2/departures/", `{"status":400,"message":"Invalid stop ID. Use https://stoppested.entur.org/ to find stop IDs."}`, 400},
		{"/api/v2/departures/60890", fmt.Sprintf(`{"url":"%s/api/v2/departures/60890","departures":[{"line":"11","scheduledDepartureTime":"2021-08-11T23:33:09.000","destination":"Risvollan via sentrum","isRealtimeData":true,"isGoingTowardsCentrum":false},{"line":"3","scheduledDepartureTime":"2021-08-11T23:38:01.000","destination":"Hallset","isRealtimeData":true,"isGoingTowardsCentrum":true}]}`, httpSrv.URL), 200},
		{"/api/v2/departures/60890?direction=inbound", fmt.Sprintf(`{"url":"%s/api/v2/departures/60890","departures":[{"line":"3","scheduledDepartureTime":"2021-08-11T23:38:01.000","destination":"Hallset","isRealtimeData":true,"isGoingTowardsCentrum":true}]}`, httpSrv.URL), 200},
	}
	for _, tt := range tests {
		data, contentType, status, err := httpGet(httpSrv.URL + tt.url)
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
		{&http.Request{Host: "baz", Header: map[string][]string{"X-Forwarded-Proto": {"https"}}}, "https://baz"},
		{&http.Request{Host: "qux", Header: map[string][]string{"X-Forwarded-Proto": {}}}, "http://qux"},
	}
	for _, tt := range tests {
		prefix := urlPrefix(tt.in)
		if prefix != tt.out {
			t.Errorf("want %s, got %s", tt.out, prefix)
		}
	}
}

const enturResponse = `{
  "data": {
    "stopPlace": {
      "id": "NSR:StopPlace:60890",
      "name": "Ila",
      "estimatedCalls": [
        {
          "realtime": true,
          "expectedDepartureTime": "2021-08-11T23:33:09+02:00",
          "actualDepartureTime": null,
          "destinationDisplay": {
            "frontText": "Risvollan via sentrum"
          },
          "serviceJourney": {
            "operator": {
              "id": "ATB:Operator:171"
            },
            "journeyPattern": {
              "directionType": "outbound",
              "line": {
                "publicCode": "11"
              }
            }
          }
        },
        {
          "realtime": true,
          "expectedDepartureTime": "2021-08-11T23:38:01+02:00",
          "actualDepartureTime": null,
          "destinationDisplay": {
            "frontText": "Hallset"
          },
          "serviceJourney": {
            "operator": {
              "id": "ATB:Operator:171"
            },
            "journeyPattern": {
              "directionType": "inbound",
              "line": {
                "publicCode": "3"
              }
            }
          }
        }
      ]
    }
  }
}`
