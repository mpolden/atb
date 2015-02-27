package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/martinp/atbapi/atb"
	"github.com/pmylund/go-cache"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Api struct {
	Client        atb.Client
	busStopsCache *cache.Cache
}

func indentJSON(req *http.Request) bool {
	_, ok := req.URL.Query()["pretty"]
	return ok
}

func marshalJSON(data interface{}, indent bool) ([]byte, error) {
	if indent {
		return json.MarshalIndent(data, "", "  ")
	}
	return json.Marshal(data)
}

func (a *Api) getBusStops() (BusStops, error) {
	const cacheKey = "stops"
	cached, ok := a.busStopsCache.Get(cacheKey)
	if ok {
		cachedBusStops, ok := cached.(BusStops)
		if !ok {
			return BusStops{}, fmt.Errorf(
				"type assertion of cached value failed")
		}
		return cachedBusStops, nil
	}
	atbBusStops, err := a.Client.GetBusStops()
	if err != nil {
		return BusStops{}, err
	}
	busStops, err := convertBusStops(atbBusStops)
	if err != nil {
		return BusStops{}, err
	}
	log.Print("Adding bus stops to cache")
	a.busStopsCache.Set(cacheKey, busStops, cache.DefaultExpiration)
	return busStops, nil
}

func (a *Api) BusStopsHandler(w http.ResponseWriter, req *http.Request) {
	busStops, err := a.getBusStops()
	if err != nil {
		// XXX: Return JSON
		http.Error(w, "Failed to get bus stops",
			http.StatusInternalServerError)
		log.Print(err)
		return
	}
	indent := indentJSON(req)
	jsonBlob, err := marshalJSON(busStops, indent)
	if err != nil {
		http.Error(w, "Failed to marshal bus stops",
			http.StatusInternalServerError)
		log.Print(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBlob)
}

func (a *Api) ForecastHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	nodeId, err := strconv.Atoi(vars["nodeId"])
	if err != nil {
		http.Error(w, "Missing or invalid nodeId", http.StatusBadRequest)
		log.Print(err)
		return
	}
	forecasts, err := a.Client.GetRealTimeForecast(nodeId)
	if err != nil {
		http.Error(w, "Failed to get forecast from upstream",
			http.StatusInternalServerError)
		log.Print(err)
		return
	}
	departures, err := convertForecasts(forecasts)
	if err != nil {
		http.Error(w, "Failed to convert forecast",
			http.StatusInternalServerError)
		log.Print(err)
		return
	}
	indent := indentJSON(req)
	jsonBlob, err := marshalJSON(departures, indent)
	if err != nil {
		http.Error(w, "Failed to marshal departures",
			http.StatusInternalServerError)
		log.Print(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBlob)
}

func New(client atb.Client) Api {
	// Cache bus stops for 1 week, check for expiration every day
	busStopsCache := cache.New(168*time.Hour, 24*time.Hour)
	return Api{
		Client:        client,
		busStopsCache: busStopsCache,
	}
}

func ListenAndServe(client atb.Client, addr string) error {
	api := New(client)
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/busstops", api.BusStopsHandler)
	r.HandleFunc("/api/v1/departures/{nodeId:[0-9]+}", api.ForecastHandler)
	http.Handle("/", r)
	return http.ListenAndServe(addr, nil)
}
