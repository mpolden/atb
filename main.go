package main

import (
	"github.com/jessevdk/go-flags"
	"github.com/martinp/atbapi/api"
	"github.com/martinp/atbapi/atb"
	"log"
	"os"
	"time"
)

func main() {
	var opts struct {
		Listen          string        `short:"l" long:"listen" description:"Listen address" value-name:"ADDRESS" default:":8080"`
		Config          string        `short:"c" long:"config" description:"Path to config file" value-name:"FILE" default:"config.json"`
		CacheStops      time.Duration `short:"s" long:"cache-stops" description:"Bus stops cache duration" value-name:"DURATION" default:"168h"`
		CacheDepartures time.Duration `short:"d" long:"cache-departures" description:"Departures cache duration" value-name:"DURATION" default:"1m"`
	}
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		os.Exit(1)
	}

	client, err := atb.NewFromConfig(opts.Config)
	if err != nil {
		log.Fatal(err)
	}

	api := api.New(client, opts.CacheStops, opts.CacheDepartures)

	log.Printf("Listening on %s", opts.Listen)
	if err := api.ListenAndServe(opts.Listen); err != nil {
		log.Fatal(err)
	}
}
