package api

import (
	"github.com/martinp/atbapi/atb"
	"reflect"
	"testing"
)

func TestConvertBusStop(t *testing.T) {
	stop := atb.BusStop{
		StopID:      100633,
		NodeID:      "16011376",
		Description: "Prof. Brochs gt",
		Longitude:   "1157514",
		Latitude:    9202874,
		MobileCode:  "16011376 (Prof.)",
		MobileName:  "Prof. (16011376)",
	}
	expected := BusStop{
		StopID:      100633,
		NodeID:      16011376,
		Description: "Prof. Brochs gt",
		Longitude:   10.398125177823237,
		Latitude:    63.4155348940887,
		MobileCode:  "16011376 (Prof.)",
		MobileName:  "Prof. (16011376)",
	}
	actual, err := convertBusStop(stop)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %+v, got %+v", expected, actual)
	}
}

func TestConvertBusStops(t *testing.T) {
	stops := atb.BusStops{
		Stops: []atb.BusStop{atb.BusStop{
			NodeID:    "16011376",
			Longitude: "1157514",
			Latitude:  9202874,
		}}}
	expected := BusStops{
		Stops: []BusStop{BusStop{
			NodeID:    16011376,
			Longitude: 10.398125177823237,
			Latitude:  63.4155348940887,
		}}}
	actual, err := convertBusStops(stops)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %+v, got %+v", expected, actual)
	}
}

func TestConvertTime(t *testing.T) {
	time, err := ConvertTime("26.02.2015 18:38")
	if err != nil {
		t.Fatal(err)
	}
	expected := "2015-02-26T18:38:00.000"
	if time != expected {
		t.Fatalf("Expected %s, got %s", expected, time)
	}
}

func TestIsRealtime(t *testing.T) {
	if !IsRealtime("prev") {
		t.Fatal("Expected true")
	}
	if !IsRealtime("Prev") {
		t.Fatal("Expected true")
	}
	if IsRealtime("foo") {
		t.Fatal("Expected false")
	}
}

func TestConvertForecast(t *testing.T) {
	forecast := atb.Forecast{
		LineID:                  "6",
		LineDescription:         "6",
		RegisteredDepartureTime: "26.02.2015 18:38",
		ScheduledDepartureTime:  "26.02.2015 18:01",
		StationForecast:         "Prev",
		Destination:             "Munkegata M5",
	}
	expected := Departure{
		LineID:                  "6",
		Destination:             "Munkegata M5",
		RegisteredDepartureTime: "2015-02-26T18:38:00.000",
		ScheduledDepartureTime:  "2015-02-26T18:01:00.000",
		IsRealtimeData:          true,
	}
	actual, err := convertForecast(forecast)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %+v, got %+v", expected, actual)
	}
}

func TestIsTowardsCentrum(t *testing.T) {
	if !IsTowardsCentrum(16011376) {
		t.Fatal("Expected true")
	}
	if IsTowardsCentrum(16010376) {
		t.Fatal("Expected false")
	}
}

func TestConvertForecasts(t *testing.T) {
	forecasts := atb.Forecasts{
		Nodes: []atb.NodeInfo{atb.NodeInfo{NodeID: "16011376"}},
		Forecasts: []atb.Forecast{atb.Forecast{
			RegisteredDepartureTime: "26.02.2015 18:38",
			ScheduledDepartureTime:  "26.02.2015 18:01",
		}}}
	expected := Departures{TowardsCentrum: true,
		Departures: []Departure{Departure{
			RegisteredDepartureTime: "2015-02-26T18:38:00.000",
			ScheduledDepartureTime:  "2015-02-26T18:01:00.000",
			IsRealtimeData:          false,
		}}}
	actual, err := convertForecasts(forecasts)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %+v, got %+v", expected, actual)
	}
}

func TestConvertCoordinates(t *testing.T) {
	// Prof. Brochs gate
	latitude, longitude := 9202565, 1157522
	lat, lon := ConvertCoordinates(latitude, longitude)
	if expected := 63.41429265308724; lat != expected {
		t.Fatalf("Expected %f, got %f", expected, lat)
	}
	if expected := 10.398197043045966; lon != expected {
		t.Fatalf("Expected %f, got %f", expected, lon)
	}

	// Ilsvika
	latitude, longitude = 9206756, 1152920
	lat, lon = ConvertCoordinates(latitude, longitude)
	if expected := 63.43113671582598; lat != expected {
		t.Fatalf("Expected %f, got %f", expected, lat)
	}
	if expected := 10.356856573670786; lon != expected {
		t.Fatalf("Expected %f, got %f", expected, lon)
	}
}
