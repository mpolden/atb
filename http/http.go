package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"time"

	"github.com/mpolden/atb/atb"
	"github.com/mpolden/atb/cache"
)

// Server represents an Server server.
type Server struct {
	Client atb.Client
	CORS   bool
	cache  *cache.Cache
	ttl
}

type ttl struct {
	departures time.Duration
	stops      time.Duration
}

func urlPrefix(r *http.Request) string {
	host := r.Host
	if host == "" {
		host = r.RemoteAddr
	}
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	forwardedProto := r.Header["X-Forwarded-Proto"]
	if len(forwardedProto) > 0 {
		scheme = forwardedProto[0]
	}
	url := url.URL{Scheme: scheme, Host: host}
	return url.String()
}

func (s *Server) getBusStops(urlPrefix string) (BusStops, bool, error) {
	const cacheKey = "stops"
	cached, hit := s.cache.Get(cacheKey)
	if hit {
		return cached.(BusStops), hit, nil
	}
	atbBusStops, err := s.Client.BusStops()
	if err != nil {
		return BusStops{}, hit, err
	}
	busStops, err := convertBusStops(atbBusStops)
	if err != nil {
		return BusStops{}, hit, err
	}
	for i := range busStops.Stops {
		busStops.Stops[i].URL = fmt.Sprintf("%s/api/v1/busstops/%d", urlPrefix, busStops.Stops[i].NodeID)
	}
	// Create a map of nodeIds
	busStops.nodeIDs = make(map[int]*BusStop, len(busStops.Stops))
	for i, s := range busStops.Stops {
		// Store a pointer to the BusStop struct
		busStops.nodeIDs[s.NodeID] = &busStops.Stops[i]
	}
	s.cache.Set(cacheKey, busStops, s.ttl.stops)
	return busStops, hit, nil
}

func (s *Server) getDepartures(urlPrefix string, nodeID int) (Departures, bool, error) {
	cacheKey := strconv.Itoa(nodeID)
	cached, hit := s.cache.Get(cacheKey)
	if hit {
		return cached.(Departures), hit, nil
	}
	forecasts, err := s.Client.Forecasts(nodeID)
	if err != nil {
		return Departures{}, hit, err
	}
	departures, err := convertForecasts(forecasts)
	if err != nil {
		return Departures{}, hit, err
	}
	departures.URL = fmt.Sprintf("%s/api/v1/departures/%d", urlPrefix, nodeID)
	s.cache.Set(cacheKey, departures, s.ttl.departures)
	return departures, hit, nil
}

func (s *Server) setCacheHeader(w http.ResponseWriter, hit bool) {
	v := "MISS"
	if hit {
		v = "HIT"
	}
	w.Header().Set("X-Cache", v)
}

// BusStopsHandler is a handler for retrieving bus stops.
func (s *Server) BusStopsHandler(w http.ResponseWriter, r *http.Request) (interface{}, *Error) {
	busStops, hit, err := s.getBusStops(urlPrefix(r))
	if err != nil {
		return nil, &Error{
			err:     err,
			Status:  http.StatusInternalServerError,
			Message: "Failed to get bus stops from AtB",
		}
	}
	s.setCacheHeader(w, hit)
	_, geojson := r.URL.Query()["geojson"]
	if geojson {
		return busStops.GeoJSON(), nil
	}
	return busStops, nil
}

// BusStopHandler is a handler for retrieving info about a bus stop.
func (s *Server) BusStopHandler(w http.ResponseWriter, r *http.Request) (interface{}, *Error) {
	nodeID, err := strconv.Atoi(filepath.Base(r.URL.Path))
	if err != nil {
		return nil, &Error{
			err:     err,
			Status:  http.StatusBadRequest,
			Message: "Invalid nodeID",
		}
	}
	busStops, hit, err := s.getBusStops(urlPrefix(r))
	if err != nil {
		return nil, &Error{
			err:     err,
			Status:  http.StatusInternalServerError,
			Message: "Failed to get bus stops from AtB",
		}
	}
	busStop, ok := busStops.nodeIDs[nodeID]
	if !ok {
		return nil, &Error{
			Status:  http.StatusNotFound,
			Message: "Unknown bus stop",
		}
	}
	s.setCacheHeader(w, hit)
	_, geojson := r.URL.Query()["geojson"]
	if geojson {
		return busStop.GeoJSON(), nil
	}
	return busStop, nil
}

// DepartureHandler is a handler for retrieving departures for a given bus stop.
func (s *Server) DepartureHandler(w http.ResponseWriter, r *http.Request) (interface{}, *Error) {
	nodeID, err := strconv.Atoi(filepath.Base(r.URL.Path))
	if err != nil {
		return nil, &Error{
			err:     err,
			Status:  http.StatusBadRequest,
			Message: "Invalid nodeID",
		}
	}
	busStops, _, err := s.getBusStops(urlPrefix(r))
	if err != nil {
		return nil, &Error{
			err:     err,
			Status:  http.StatusInternalServerError,
			Message: "Failed to get bus stops from AtB",
		}
	}
	_, ok := busStops.nodeIDs[nodeID]
	if !ok {
		return nil, &Error{
			Status:  http.StatusNotFound,
			Message: "Unknown bus stop",
		}
	}
	departures, hit, err := s.getDepartures(urlPrefix(r), nodeID)
	if err != nil {
		return nil, &Error{
			err:     err,
			Status:  http.StatusInternalServerError,
			Message: "Failed to get departures from AtB",
		}
	}
	s.setCacheHeader(w, hit)
	return departures, nil
}

// DeparturesHandler lists all known departures.
func (s *Server) DeparturesHandler(w http.ResponseWriter, r *http.Request) (interface{}, *Error) {
	busStops, hit, err := s.getBusStops(urlPrefix(r))
	if err != nil {
		return nil, &Error{
			err:     err,
			Status:  http.StatusInternalServerError,
			Message: "Failed to get bus stops from AtB",
		}
	}
	s.setCacheHeader(w, hit)
	var urls struct {
		URLs []string `json:"urls"`
	}
	urls.URLs = make([]string, len(busStops.Stops))
	for i, stop := range busStops.Stops {
		urls.URLs[i] = fmt.Sprintf("%s/api/v1/departures/%d", urlPrefix(r), stop.NodeID)
	}
	return urls, nil
}

// DefaultHandler lists known URLs.
func (s *Server) DefaultHandler(w http.ResponseWriter, r *http.Request) (interface{}, *Error) {
	if r.URL.Path != "/" {
		return nil, &Error{Status: http.StatusNotFound, Message: "Resource not found"}
	}
	prefix := urlPrefix(r)
	busStopsURL := fmt.Sprintf("%s/api/v1/busstops", prefix)
	departuresURL := fmt.Sprintf("%s/api/v1/departures", prefix)
	return struct {
		URLs []string `json:"urls"`
	}{
		[]string{busStopsURL, departuresURL},
	}, nil
}

// New returns a new Server using client to communicate with AtB. stopTTL and departureTTL control the cache TTL bus
// stops and departures.
func New(client atb.Client, stopTTL, departureTTL time.Duration, cors bool) Server {
	cache := cache.New(time.Minute)
	return Server{
		Client: client,
		CORS:   cors,
		cache:  cache,
		ttl: ttl{
			stops:      stopTTL,
			departures: departureTTL,
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
		out, err := json.Marshal(e)
		if err != nil {
			// Should never happen
			panic(err)
		}
		w.WriteHeader(e.Status)
		w.Write(out)
	} else {
		out, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		w.Write(out)
	}
}

func requestFilter(next http.Handler, cors bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if cors {
			w.Header().Set("Access-Control-Allow-Methods", "GET")
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		next.ServeHTTP(w, r)
	})
}

// Handler returns a root handler for the API.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/api/v1/busstops", appHandler(s.BusStopsHandler))
	mux.Handle("/api/v1/busstops/", appHandler(s.BusStopHandler))
	mux.Handle("/api/v1/departures", appHandler(s.DeparturesHandler))
	mux.Handle("/api/v1/departures/", appHandler(s.DepartureHandler))
	mux.Handle("/", appHandler(s.DefaultHandler))
	return requestFilter(mux, s.CORS)
}

// ListenAndServe listens on the TCP network address addr and serves the API.
func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s.Handler())
}
