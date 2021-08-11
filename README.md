# atb

![Build Status](https://github.com/mpolden/atb/workflows/ci/badge.svg)

A minimal API for bus data in Trondheim, Norway. This API proxies requests to
AtB/Entur APIs and converts the responses into a sane JSON format.

Responses from the proxied APIs are cached. By default bus stops will be cached
for 1 week and departures for 1 minute.

As of mid-August 2021 the SOAP-based AtB API appears to no longer return any
departures. According to [this blog post on open
data](https://beta.atb.no/blogg/apne-data-og-atb) it appears AtB now provides
data through [Entur](https://developer.entur.org/).

Version 1 of this API will remain implemented for now, but likely won't return
any usable data.

Version 2 has been implemented and proxies requests to Entur. Version 2 differs
from version 1 in the following ways:

* There is no version 2 variant of `/api/v1/busstops`. Use
  https://stoppested.entur.org/ to find valid stop IDs.
* Node/stop IDs have changed so the old ones (e.g. `16011376`) cannot be used in
  version 2.
* The `registeredDepartureTime` field may be omitted.
* The `isGoingTowardsCentrum` field has moved to the departure object.

Both version 1 and 2 of this API aims to be compatible with
[BusBuddy](https://github.com/norrs/busbuddy) (which appears to be defunct).

## Usage

```
$ atb -h
Usage of atb:
  -c string
    	Path to config file (default "config.json")
  -d string
    	Departure cache duration (default "1m")
  -l string
    	Listen address (default ":8080")
  -s string
    	Bus stop cache duration (default "168h")
  -x	Allow requests from other domains
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
    "https://mpolden.no/atb/v1/departures",
    "https://mpolden.no/atb/v2/departures"
  ]
}
```

### `/api/v2/departures`

List departures from the given bus stop, identified by a stop ID. Use
https://stoppested.entur.org to find stop IDs, for example `43501` (the number
part of `NSR:StopPlace:43501`) for Dronningens gate.

Departures in any direction are included by default. Add the parameter
`direction=inbound` or `direction=outbound` to filter departures towards, or
away from, the city centre.

```
$ curl https://mpolden.no/atb/v2/departures/41613 | jq .

{
  "url": "https://mpolden.no/atb/v2/departures/41613",
  "departures": [
    {
      "line": "71",
      "scheduledDepartureTime": "2021-08-11T23:49:38.000",
      "destination": "Dora",
      "isRealtimeData": true,
      "isGoingTowardsCentrum": true
    },
    ...
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
