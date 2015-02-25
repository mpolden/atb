package main

import (
	"encoding/json"
	"github.com/jessevdk/go-flags"
	"io/ioutil"
	"log"
	"os"
)

func ReadConfig(name string) (Atb, error) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return Atb{}, err
	}
	var atb Atb
	if err := json.Unmarshal(data, &atb); err != nil {
		return Atb{}, err
	}
	return atb, nil
}

func main() {
	var opts struct {
		Config    string `short:"c" long:"config" description:"Path to config file" value-name:"FILE" default:"config.json"`
		Templates string `short:"t" long:"templates" description:"Path to request templates directory" value-name:"PATH" default:"templates"`
	}
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		os.Exit(1)
	}

	a, err := ReadConfig(opts.Config)
	if err != nil {
		log.Fatal(err)
	}
	methods, err := NewMethods(opts.Templates)
	if err != nil {
		log.Fatal(err)
	}
	a.Methods = methods
	_, err = a.GetBusStops()
	if err != nil {
		log.Fatal(err)
	}
}
