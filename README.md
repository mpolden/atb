# atbapi

![Build Status](https://github.com/mpolden/atbapi/workflows/ci/badge.svg)

A minimal API for bus data in Trondheim, Norway. This API proxies requests to
the AtB public API and converts the responses into a sane JSON format.

Responses from AtBs public API are cached. By default bus stops will be cached
for 1 week and departures for 1 minute.

If you want to setup this service yourself, you need to request access to the
[SOAP API provided by
AtB](https://www.atb.no/sanntid/apne-data-article5852-1381.html) (Norwegian).
When granted access, you'll receive a username and password (see config example
below).

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

## API

### `/`

Lists all available API routes.

Example:

```
$ curl https://mpolden.no/atb/ | jq .
{
  "urls": [
    "https://mpolden.no/atb/v1/busstops",
    "https://mpolden.no/atb/v1/departures"
  ]
}
```


### `/api/v1/busstops`

Lists all known bus stops.

Example:

```
$ curl https://mpolden.no/atb/v1/busstops | jq .
{
  "stops": [
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

### `/api/v1/busstops/{node-id}`

Information about the given bus stop, identified by a node ID.

Example:

```
$ curl https://mpolden.no/atb/v1/busstops/16011376 | jq .
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
$ curl https://mpolden.no/atb/v1/busstops/16011376?geojson | jq .
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

### `/api/v1/departures`

Lists departure URLs for all known bus stops.

Example:

```
$ curl -s https://mpolden.no/atb/v1/departures | jq .
{
  "urls": [
    "https://mpolden.no/atb/v1/departures/15057011",
    ...
  ]
}
```

### `/api/v1/departures/{node-id}`

Lists all departures for the given bus stop, identified by a node ID.

Example:

```
$ curl https://mpolden.no/atb/v1/departures/16011376 | jq .
{
  "isGoingTowardsCentrum": true,
  "departures": [
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
