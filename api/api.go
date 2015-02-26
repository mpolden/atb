package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/martinp/atbapi/atb"
	"net/http"
	"strconv"
)

type Api struct {
	Client atb.Client
}

func (a *Api) BusStopsHandler(w http.ResponseWriter, req *http.Request) {
	_busStops, err := a.Client.GetBusStops()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	busStops, err := convertBusStops(_busStops)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonBlob, err := json.Marshal(busStops)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBlob)
}

func (a *Api) ForecastHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	nodeId, err := strconv.Atoi(vars["nodeId"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}
	forecasts, err := a.Client.GetRealTimeForecast(nodeId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	departures, err := convertForecasts(forecasts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonBlob, err := json.Marshal(departures)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBlob)
}

func New(client atb.Client) Api {
	api := Api{Client: client}
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/busstops", api.BusStopsHandler)
	r.HandleFunc("/api/v1/departures/{nodeId:[0-9]+}", api.ForecastHandler)
	http.Handle("/", r)
	return api
}

func (a *Api) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, nil)
}