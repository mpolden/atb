# atbapi

[![Build Status](https://travis-ci.org/mpolden/atbapi.svg)](https://travis-ci.org/mpolden/atbapi)

A minimal API for bus data in Trondheim, Norway. This API proxies requests to
the AtB public API and converts the responses into a sane JSON format.

Responses from AtBs public API are cached. By default bus stops will be cached
for 1 week and departures for 1 minute.

If you want to setup this service yourself, you need to request access to the
SOAP API provided by AtB [here](https://www.atb.no/aapne-data/category419.html)
(Norwegian). When granted access, you'll receive a username and password (see
config example below).

The API aims to be compatible with [BusBuddy](https://github.com/norrs/busbuddy)
(which appears to be defunct).

## Usage

```
$ atbapi -h
Usage:
  atbapi [OPTIONS]

Application Options:
  -l, --listen=ADDRESS               Listen address (:8080)
  -c, --config=FILE                  Path to config file (config.json)
  -s, --cache-stops=DURATION         Bus stops cache duration (168h)
  -d, --cache-departures=DURATION    Departures cache duration (1m)
  -x, --cors                         Allow requests from other domains (false)

Help Options:
  -h, --help                         Show this help message
```

## Example config

```
{
  "Username": "username",
  "Password": "password"
}
```

## API usage

The route `/api/v1/busstops` returns a list of all known bus stops.

Example:

```
$ curl 'https://atbapi.tar.io/api/v1/busstops' | jq .
{
  "stops": [
    ...
    {
      "stopId": 100633,
      "nodeId": 16011376,
      "description": "Prof. Brochs gt",
      "longitude": 10.398125177823237,
      "latitude": 63.4155348940887,
      "mobileCode": "16011376 (Prof.)",
      "mobileName": "Prof. (16011376)"
    },
    ...
  ]
}
```

The route `/api/v1/busstops/<nodeID>` returns information about a given bus
stop, identified by a node ID.

Example:

```
$ curl 'https://atbapi.tar.io/api/v1/busstops/16011376' | jq .
{
  "stopId": 100633,
  "nodeId": 16011376,
  "description": "Prof. Brochs gt",
  "longitude": 10.398126,
  "latitude": 63.415535,
  "mobileCode": "16011376 (Prof.)",
  "mobileName": "Prof. (16011376)"
}
```

As [GeoJSON](http://geojson.org/):

```
$ curl 'https://atbapi.tar.io/api/v1/busstops/16011376?geojson' | jq .
{
  "type": "Feature",
  "geometry": {
    "type": "Point",
    "coordinates": [
      10.398126,
      63.415535
    ]
  },
  "properties": {
    "busstop": {
      "stopId": 100633,
      "nodeId": 16011376,
      "description": "Prof. Brochs gt",
      "longitude": 10.398126,
      "latitude": 63.415535,
      "mobileCode": "16011376 (Prof.)",
      "mobileName": "Prof. (16011376)"
    },
    "name": "Prof. Brochs gt"
  }
}
```

The route `/api/v1/departures/<nodeID>` returns a list of departures for a given bus
stop, identified by a node ID.

Example:

```
$ curl 'https://atbapi.tar.io/api/v1/departures/16011376' | jq .
{
  "isGoingTowardsCentrum": true,
  "departures": [
    ...
    {
      "line": "36",
      "registeredDepartureTime": "2015-02-26T22:55:00.000",
      "scheduledDepartureTime": "2015-02-26T22:54:00.000",
      "destination": "Munkegata M4",
      "isRealtimeData": true
    },
    ...
  ]
}
```
