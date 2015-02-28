# atbapi

[![Build Status](https://travis-ci.org/martinp/atbapi.png)](https://travis-ci.org/martinp/atbapi)

A minimal JSON API for bus data in Trondheim, Norway. This API proxies requests
to the AtB public API and converts the responses into a sane format.

Responses from AtBs public API are cached. By default bus stops will be cached
for 1 week and departures for 1 minute.

You can request access to the SOAP API provided by AtB
[here](https://www.atb.no/aapne-data/category419.html) (Norwegian).

The API aims to be compatible with [BusBuddy](https://github.com/norrs/busbuddy).

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

The route `/api/v1/busstops` returns a list of all known bus stops. All routes
support the parameter `?pretty` to return pretty-printed JSON.

Example:

```
$ curl 'http://localhost:8080/api/v1/busstops?pretty'
{
  "stops": [
    ...
    {
      "stopId": 100633,
      "nodeId": 16011376,
      "description": "Prof. Brochs gt",
      "longitude": 1157514,
      "latitude": 9202874,
      "mobileCode": "16011376 (Prof.)",
      "mobileName": "Prof. (16011376)"
    },
    ...
  ]
}
```

The route `/api/v1/departures` returns a list of departures for a given bus
stop, identified by `nodeId`.

Example:

```
$ curl 'http://localhost:8080/api/v1/departures/16011376?pretty'
{
  "isGoingTowardsCentrum": true,
  "departures": [
    {
      "line": "8",
      "registeredDepartureTime": "2015-02-26T22:52:00.000",
      "scheduledDepartureTime": "2015-02-26T22:53:00.000",
      "destination": "Kongens gate K1",
      "isRealtimeData": true
    },
    {
      "line": "36",
      "registeredDepartureTime": "2015-02-26T22:55:00.000",
      "scheduledDepartureTime": "2015-02-26T22:54:00.000",
      "destination": "Munkegata M4",
      "isRealtimeData": true
    },
    {
      "line": "38",
      "registeredDepartureTime": "2015-02-26T23:01:00.000",
      "scheduledDepartureTime": "2015-02-26T22:57:00.000",
      "destination": "Dronningens gt.",
      "isRealtimeData": true
    },
    {
      "line": "6",
      "registeredDepartureTime": "2015-02-26T23:02:00.000",
      "scheduledDepartureTime": "2015-02-26T23:01:00.000",
      "destination": "Munkegata M5",
      "isRealtimeData": true
    },
    {
      "line": "46",
      "registeredDepartureTime": "2015-02-26T23:08:00.000",
      "scheduledDepartureTime": "2015-02-26T23:05:00.000",
      "destination": "Pirbadet",
      "isRealtimeData": true
    }
  ]
}
```
