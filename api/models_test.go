package api

import (
	"github.com/martinp/atbapi/atb"
	"reflect"
	"testing"
)

func TestConvertBusStop(t *testing.T) {
	stop := atb.BusStop{
		StopId:      100633,
		NodeId:      "16011376",
		Description: "Prof. Brochs gt",
		Longitude:   "1157514",
		Latitude:    9202874,
		MobileCode:  "16011376 (Prof.)",
		MobileName:  "Prof. (16011376)",
	}
	expected := BusStop{
		StopId:      100633,
		NodeId:      16011376,
		Description: "Prof. Brochs gt",
		Longitude:   1157514,
		Latitude:    9202874,
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
			NodeId:    "16011376",
			Longitude: "1157514",
			Latitude:  9202874,
		}}}
	expected := BusStops{
		Stops: []BusStop{BusStop{
			NodeId:    16011376,
			Longitude: 1157514,
			Latitude:  9202874,
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
	time, err := convertTime("26.02.2015 18:38")
	if err != nil {
		t.Fatal(err)
	}
	expected := "2015-02-26T18:38:00.000"
	if time != expected {
		t.Fatalf("Expected %s, got %s")
	}
}

func TestIsRealtime(t *testing.T) {
	if !isRealtime("prev") {
		t.Fatal("Expected true")
	}
	if !isRealtime("Prev") {
		t.Fatal("Expected true")
	}
	if isRealtime("foo") {
		t.Fatal("Expected false")
	}
}

func TestConvertForecast(t *testing.T) {
	forecast := atb.Forecast{
		LineId:                  "6",
		LineDescription:         "6",
		RegisteredDepartureTime: "26.02.2015 18:38",
		ScheduledDepartureTime:  "26.02.2015 18:01",
		StationForecast:         "Prev",
		Destination:             "Munkegata M5",
	}
	expected := Departure{
		LineId:                  "6",
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
	if !isTowardsCentrum(16011376) {
		t.Fatal("Expected true")
	}
	if isTowardsCentrum(16010376) {
		t.Fatal("Expected false")
	}
}

func TestConvertForecasts(t *testing.T) {
	forecasts := atb.Forecasts{
		Nodes: []atb.NodeInfo{atb.NodeInfo{NodeId: "16011376"}},
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
