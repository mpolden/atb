package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/martinp/atbapi/atb"
	"log"
	"net/http"
	"strconv"
)

type Api struct {
	Client atb.Client
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

func (a *Api) BusStopsHandler(w http.ResponseWriter, req *http.Request) {
	_busStops, err := a.Client.GetBusStops()
	if err != nil {
		http.Error(w, "Failed to get bus stops from upstream",
			http.StatusInternalServerError)
		log.Print(err)
		return
	}
	busStops, err := convertBusStops(_busStops)
	if err != nil {
		http.Error(w, "Failed to convert bus stops",
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

func ListenAndServe(client atb.Client, addr string) error {
	api := Api{Client: client}
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/busstops", api.BusStopsHandler)
	r.HandleFunc("/api/v1/departures/{nodeId:[0-9]+}", api.ForecastHandler)
	http.Handle("/", r)
	return http.ListenAndServe(addr, nil)
}
