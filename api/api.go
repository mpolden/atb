package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/martinp/atbapi/atb"
	cache "github.com/pmylund/go-cache"
)

// An API defines parameters for running an API server.
type API struct {
	Client atb.Client
	CORS   bool
	cache  *cache.Cache
	expiration
}

type expiration struct {
	departures time.Duration
	stops      time.Duration
}

func marshal(data interface{}, indent bool) ([]byte, error) {
	if indent {
		return json.MarshalIndent(data, "", "  ")
	}
	return json.Marshal(data)
}

func (a *API) getBusStops() (BusStops, bool, error) {
	const cacheKey = "stops"
	cached, hit := a.cache.Get(cacheKey)
	if hit {
		cachedBusStops, ok := cached.(BusStops)
		if !ok {
			return BusStops{}, false, fmt.Errorf(
				"type assertion of cached value failed")
		}
		return cachedBusStops, hit, nil
	}
	atbBusStops, err := a.Client.GetBusStops()
	if err != nil {
		return BusStops{}, hit, err
	}
	busStops, err := convertBusStops(atbBusStops)
	if err != nil {
		return BusStops{}, hit, err
	}
	// Create a map of nodeIds
	busStops.nodeIDs = make(map[int]*BusStop, len(busStops.Stops))
	for i, s := range busStops.Stops {
		// Store a pointer to the BusStop struct
		busStops.nodeIDs[s.NodeID] = &busStops.Stops[i]
	}
	a.cache.Set(cacheKey, busStops, a.expiration.stops)
	return busStops, hit, nil
}

func (a *API) getDepartures(nodeID int) (Departures, bool, error) {
	cacheKey := strconv.Itoa(nodeID)
	cached, hit := a.cache.Get(cacheKey)
	if hit {
		cachedDepartures, ok := cached.(Departures)
		if !ok {
			return Departures{}, false, fmt.Errorf(
				"type assertion of cached value failed")
		}
		return cachedDepartures, hit, nil
	}
	forecasts, err := a.Client.GetRealTimeForecast(nodeID)
	if err != nil {
		return Departures{}, hit, err
	}
	departures, err := convertForecasts(forecasts)
	if err != nil {
		return Departures{}, hit, err
	}
	a.cache.Set(cacheKey, departures, cache.DefaultExpiration)
	return departures, hit, nil
}

func (a *API) setCacheHeader(w http.ResponseWriter, hit bool) {
	v := "MISS"
	if hit {
		v = "HIT"
	}
	w.Header().Set("X-Cache", v)
}

// BusStopsHandler is a handler for retrieving bus stops.
func (a *API) BusStopsHandler(w http.ResponseWriter, req *http.Request) (interface{}, *Error) {
	busStops, hit, err := a.getBusStops()
	if err != nil {
		return nil, &Error{
			err:     err,
			Status:  http.StatusInternalServerError,
			Message: "failed to get bus stops from atb",
		}
	}
	a.setCacheHeader(w, hit)
	_, geojson := req.URL.Query()["geojson"]
	if geojson {
		return busStops.GeoJSON(), nil
	}

	return busStops, nil
}

// BusStopHandler is a handler for retrieving info about a bus stop.
func (a *API) BusStopHandler(w http.ResponseWriter, req *http.Request) (interface{}, *Error) {
	vars := mux.Vars(req)
	nodeID, err := strconv.Atoi(vars["nodeID"])
	if err != nil {
		return nil, &Error{
			err:     err,
			Status:  http.StatusBadRequest,
			Message: "missing or invalid nodeID",
		}
	}
	busStops, hit, err := a.getBusStops()
	if err != nil {
		return nil, &Error{
			err:     err,
			Status:  http.StatusInternalServerError,
			Message: "failed to get bus stops from atb",
		}
	}
	busStop, ok := busStops.nodeIDs[nodeID]
	if !ok {
		msg := fmt.Sprintf("bus stop with nodeID=%d not found", nodeID)
		return nil, &Error{
			err:     err,
			Status:  http.StatusNotFound,
			Message: msg,
		}
	}
	a.setCacheHeader(w, hit)
	_, geojson := req.URL.Query()["geojson"]
	if geojson {
		return busStop.GeoJSON(), nil
	}
	return busStop, nil
}

// DeparturesHandler is a handler for retrieving departures.
func (a *API) DeparturesHandler(w http.ResponseWriter, req *http.Request) (interface{}, *Error) {
	vars := mux.Vars(req)
	nodeID, err := strconv.Atoi(vars["nodeID"])
	if err != nil {
		return nil, &Error{
			err:     err,
			Status:  http.StatusBadRequest,
			Message: "missing or invalid nodeID",
		}
	}
	busStops, hit, err := a.getBusStops()
	if err != nil {
		return nil, &Error{
			err:     err,
			Status:  http.StatusInternalServerError,
			Message: "could not get bus stops from atb",
		}
	}
	_, knownBusStop := busStops.nodeIDs[nodeID]
	if !knownBusStop {
		msg := fmt.Sprintf("bus stop with nodeID=%d not found", nodeID)
		return nil, &Error{
			err:     err,
			Status:  http.StatusNotFound,
			Message: msg,
		}
	}
	departures, hit, err := a.getDepartures(nodeID)
	if err != nil {
		return nil, &Error{
			err:     err,
			Status:  http.StatusInternalServerError,
			Message: "could not get departures from atb",
		}
	}
	a.setCacheHeader(w, hit)
	return departures, nil
}

// NotFoundHandler handles requests to invalid routes.
func (a *API) NotFoundHandler(w http.ResponseWriter, req *http.Request) (interface{}, *Error) {
	return nil, &Error{
		err:     nil,
		Status:  http.StatusNotFound,
		Message: "route not found",
	}
}

// New returns an new API using client to communicate with AtB. stopsExpiration
// and depExpiration control the cache expiration times for bus stops and
// departures.
func New(client atb.Client, stopsExpiration, depExpiration time.Duration, cors bool) API {
	cache := cache.New(depExpiration, 30*time.Second)
	return API{
		Client: client,
		CORS:   cors,
		cache:  cache,
		expiration: expiration{
			stops:      stopsExpiration,
			departures: depExpiration,
		},
	}
}

type appHandler func(http.ResponseWriter, *http.Request) (interface{}, *Error)

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, e := fn(w, r)
	if e != nil { // e is *Error, not os.Error.
		if e.err != nil {
			log.Print(e.err)
		}
		jsonBlob, err := marshal(e, true)
		if err != nil {
			// Should never happen
			panic(err)
		}
		w.WriteHeader(e.Status)
		w.Write(jsonBlob)
	} else {
		indent := context.Get(r, "indent").(bool)
		jsonBlob, err := marshal(data, indent)
		if err != nil {
			panic(err)
		}
		w.Write(jsonBlob)
	}
}

func requestFilter(next http.Handler, cors bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, indent := r.URL.Query()["pretty"]
		context.Set(r, "indent", indent)
		w.Header().Set("Content-Type", "application/json")
		if cors {
			w.Header().Set("Access-Control-Allow-Methods", "GET")
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		next.ServeHTTP(w, r)
	})
}

// ListenAndServe listens on the TCP network address addr and starts serving the
// API.
func (a *API) ListenAndServe(addr string) error {
	r := mux.NewRouter()
	r.Handle("/api/v1/busstops", appHandler(a.BusStopsHandler))
	r.Handle("/api/v1/busstops/{nodeID:[0-9]+}", appHandler(a.BusStopHandler))
	r.Handle("/api/v1/departures/{nodeID:[0-9]+}",
		appHandler(a.DeparturesHandler))
	r.NotFoundHandler = appHandler(a.NotFoundHandler)
	http.Handle("/", requestFilter(r, a.CORS))
	return http.ListenAndServe(addr, nil)
}
