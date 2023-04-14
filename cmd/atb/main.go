package main

import (
	"flag"
	"log"
	"time"

	"github.com/mpolden/atb/entur"
	"github.com/mpolden/atb/http"
)

func init() {
	log.SetPrefix("atb: ")
	log.SetFlags(log.Lshortfile)
}

func mustParseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		log.Fatal(err)
	}
	return d
}

func main() {
	listen := flag.String("l", ":8080", "Listen address")
	stopTTL := flag.String("s", "168h", "Bus stop cache duration")
	departureTTL := flag.String("d", "1m", "Departure cache duration")
	cors := flag.Bool("x", false, "Allow requests from other domains")
	flag.Parse()

	entur := entur.New("")
	server := http.New(entur, mustParseDuration(*stopTTL), mustParseDuration(*departureTTL), *cors)

	log.Printf("Listening on %s", *listen)
	if err := server.ListenAndServe(*listen); err != nil {
		log.Fatal(err)
	}
}
