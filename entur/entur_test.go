package entur

import (
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"
)

func TestParseDepartures(t *testing.T) {
	testFile := filepath.Join("testdata", "ilsvika.json")
	json, err := ioutil.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}
	d, err := parseDepartures(json)
	if err != nil {
		t.Fatal(err)
	}
	cest := time.FixedZone("CEST", 7200)
	expected := []Departure{
		{
			Line:                    "21",
			RegisteredDepartureTime: time.Time{},
			ScheduledDepartureTime:  time.Date(2021, 8, 11, 21, 19, 0, 0, cest),
			Destination:             "Pirbadet via sentrum",
			IsRealtime:              false,
			Inbound:                 false,
		},
		{
			Line:                    "21",
			RegisteredDepartureTime: time.Time{},
			ScheduledDepartureTime:  time.Date(2021, 8, 11, 22, 19, 0, 0, cest),
			Destination:             "Pirbadet via sentrum",
			IsRealtime:              true,
			Inbound:                 false,
		},
	}
	for i := 0; i < len(expected); i++ {
		got := d[i]
		want := expected[i]
		if want.Line != got.Line {
			t.Errorf("#%d: want Line = %q, got %q", i, want.Line, got.Line)
		}
		if !want.RegisteredDepartureTime.Equal(got.RegisteredDepartureTime) {
			t.Errorf("#%d: want RegisteredDepartureTime = %q, got %q", i, want.RegisteredDepartureTime, got.RegisteredDepartureTime)
		}
		if !want.ScheduledDepartureTime.Equal(got.ScheduledDepartureTime) {
			t.Errorf("#%d: want ScheduledDepartureTime = %q, got %q", i, want.ScheduledDepartureTime, got.ScheduledDepartureTime)
		}
		if want.Destination != got.Destination {
			t.Errorf("#%d: want Destination = %q, got %q", i, want.Destination, got.Destination)
		}
		if want.IsRealtime != got.IsRealtime {
			t.Errorf("#%d: want IsRealtime = %t, got %t", i, want.IsRealtime, got.IsRealtime)
		}
		if want.Inbound != got.Inbound {
			t.Errorf("#%d: want Inbound = %t, got %t", i, want.Inbound, got.Inbound)
		}
	}
}
