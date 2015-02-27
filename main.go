package main

import (
	"github.com/jessevdk/go-flags"
	"github.com/martinp/atbapi/api"
	"github.com/martinp/atbapi/atb"
	"log"
	"os"
)

func main() {
	var opts struct {
		Listen string `short:"l" long:"listen" description:"Listen address" default:":8080"`
		Config string `short:"c" long:"config" description:"Path to config file" value-name:"FILE" default:"config.json"`
	}
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		os.Exit(1)
	}

	client, err := atb.NewFromConfig(opts.Config)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on %s", opts.Listen)
	if err := api.ListenAndServe(client, opts.Listen); err != nil {
		log.Fatal(err)
	}
}
