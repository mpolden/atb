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
	Client          atb.Client
	busStopsCache   *cache.Cache
	departuresCache *cache.Cache
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

func errorJSON(w http.ResponseWriter, message string, code int, indent bool) {
	apiError := Error{
		Message: message,
		Status:  code,
	}
	data, err := marshalJSON(apiError, indent)
	if err != nil {
		// If the error marshalling fails it's time to call it a day
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
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
	// Create a map of nodeIds
	busStops.nodeIds = make(map[int]struct{}, len(busStops.Stops))
	for _, s := range busStops.Stops {
		busStops.nodeIds[s.NodeId] = struct{}{}
	}
	a.busStopsCache.Set(cacheKey, busStops, cache.DefaultExpiration)
	return busStops, nil
}

func (a *Api) getDepartures(nodeId int) (Departures, error) {
	cacheKey := string(nodeId)
	cached, ok := a.departuresCache.Get(cacheKey)
	if ok {
		cachedDepartures, ok := cached.(Departures)
		if !ok {
			return Departures{}, fmt.Errorf(
				"type assertion of cached value failed")
		}
		return cachedDepartures, nil
	}
	forecasts, err := a.Client.GetRealTimeForecast(nodeId)
	if err != nil {
		return Departures{}, err
	}
	departures, err := convertForecasts(forecasts)
	if err != nil {
		return Departures{}, err
	}
	a.departuresCache.Set(cacheKey, departures, cache.DefaultExpiration)
	return departures, nil
}

func (a *Api) BusStopsHandler(w http.ResponseWriter, req *http.Request) {
	indent := indentJSON(req)
	busStops, err := a.getBusStops()
	if err != nil {
		errorJSON(w, "Failed to get bus stops",
			http.StatusInternalServerError, indent)
		log.Print(err)
		return
	}
	jsonBlob, err := marshalJSON(busStops, indent)
	if err != nil {
		errorJSON(w, "Failed to marshal bus stops",
			http.StatusInternalServerError, indent)
		log.Print(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBlob)
}

func (a *Api) DeparturesHandler(w http.ResponseWriter, req *http.Request) {
	indent := indentJSON(req)
	vars := mux.Vars(req)
	nodeId, err := strconv.Atoi(vars["nodeId"])
	if err != nil {
		errorJSON(w, "Missing or invalid nodeId", http.StatusBadRequest,
			indent)
		log.Print(err)
		return
	}

	busStops, err := a.getBusStops()
	if err != nil {
		errorJSON(w, "Could not get bus stops",
			http.StatusInternalServerError, indent)
		log.Print(err)
		return
	}

	_, knownBusStop := busStops.nodeIds[nodeId]
	if !knownBusStop {
		errorJSON(w, fmt.Sprintf("Bus stop with nodeId=%d not found",
			nodeId), http.StatusNotFound, indent)
		return
	}

	departures, err := a.getDepartures(nodeId)
	if err != nil {
		errorJSON(w, "Could not get departures",
			http.StatusInternalServerError, indent)
		log.Print(err)
		return
	}

	jsonBlob, err := marshalJSON(departures, indent)
	if err != nil {
		errorJSON(w, "Failed to marshal departures",
			http.StatusInternalServerError, indent)
		log.Print(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBlob)
}

func New(client atb.Client) Api {
	// Cache bus stops for 1 week, check for expiration every day
	busStopsCache := cache.New(168*time.Hour, 24*time.Hour)
	// Cache departures for 1 minute, check for expiration every 30 seconds
	departuresCache := cache.New(1*time.Minute, 30*time.Second)
	return Api{
		Client:          client,
		busStopsCache:   busStopsCache,
		departuresCache: departuresCache,
	}
}

func ListenAndServe(client atb.Client, addr string) error {
	api := New(client)
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/busstops", api.BusStopsHandler)
	r.HandleFunc("/api/v1/departures/{nodeId:[0-9]+}", api.DeparturesHandler)
	http.Handle("/", r)
	return http.ListenAndServe(addr, nil)
}
