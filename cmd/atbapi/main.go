package main

import (
	"log"
	"os"
	"time"

	flags "github.com/jessevdk/go-flags"
	"github.com/mpolden/atbapi/atb"
	"github.com/mpolden/atbapi/http"
)

func main() {
	var opts struct {
		Listen       string        `short:"l" long:"listen" description:"Listen address" value-name:"ADDRESS" default:":8080"`
		Config       string        `short:"c" long:"config" description:"Path to config file" value-name:"FILE" default:"config.json"`
		StopTTL      time.Duration `short:"s" long:"stops-ttl" description:"Bus stop cache duration" value-name:"DURATION" default:"168h"`
		DepartureTTL time.Duration `short:"d" long:"departure-ttl" description:"Departure cache duration" value-name:"DURATION" default:"1m"`
		CORS         bool          `short:"x" long:"cors" description:"Allow requests from other domains"`
	}
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		os.Exit(1)
	}

	client, err := atb.NewFromConfig(opts.Config)
	if err != nil {
		log.Fatal(err)
	}

	server := http.New(client, opts.StopTTL, opts.DepartureTTL, opts.CORS)

	log.Printf("Listening on %s", opts.Listen)
	if err := server.ListenAndServe(opts.Listen); err != nil {
		log.Fatal(err)
	}
}