package main

import (
	"log"
	"os"
	"time"

	flags "github.com/jessevdk/go-flags"
	"github.com/mpolden/atbapi/api"
	"github.com/mpolden/atbapi/atb"
)

func main() {
	var opts struct {
		Listen          string        `short:"l" long:"listen" description:"Listen address" value-name:"ADDRESS" default:":8080"`
		Config          string        `short:"c" long:"config" description:"Path to config file" value-name:"FILE" default:"config.json"`
		CacheStops      time.Duration `short:"s" long:"cache-stops" description:"Bus stops cache duration" value-name:"DURATION" default:"168h"`
		CacheDepartures time.Duration `short:"d" long:"cache-departures" description:"Departures cache duration" value-name:"DURATION" default:"1m"`
		CORS            bool          `short:"x" long:"cors" description:"Allow requests from other domains"`
	}
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		os.Exit(1)
	}

	client, err := atb.NewFromConfig(opts.Config)
	if err != nil {
		log.Fatal(err)
	}

	api := api.New(client, opts.CacheStops, opts.CacheDepartures, opts.CORS)

	log.Printf("Listening on %s", opts.Listen)
	if err := api.ListenAndServe(opts.Listen); err != nil {
		log.Fatal(err)
	}
}
