# atb

![Build Status](https://github.com/mpolden/atb/workflows/ci/badge.svg)

A minimal API for bus data in Trondheim, Norway. This API proxies requests to
Entur APIs and converts the responses into a sane JSON format.

Responses from the proxied APIs are cached. By default bus stops will be cached
for 1 week and departures for 1 minute.

As of mid-August 2021 the SOAP-based AtB API no longer returns any departure
data. According to [this blog post on open
data](https://beta.atb.no/blogg/apne-data-og-atb) it appears the preferred API
is now [Entur](https://developer.entur.org/). The `/api/v1/` paths have
therefore been removed.

Version 2 has been implemented and proxies requests to Entur instead. These are
the changes in version 2:

* There is no version 2 variant of `/api/v1/busstops`. Use
  https://stoppested.entur.org/ to find valid stop IDs.
* Entur uses different stop IDs so old ones, such as `16011376`, cannot be used
  in version 2. A stop includes departures in both directions by default so
  there is no longer a unique stop for each direction.
* The `registeredDepartureTime` field may be omitted.
* The `isGoingTowardsCentrum` field has moved to the departure object.

This API aims to be compatible with
[BusBuddy](https://github.com/norrs/busbuddy) (which appears to be defunct).

## Usage

```
$ atb -h
Usage of atb:
  -d string
    	Departure cache duration (default "1m")
  -l string
    	Listen address (default ":8080")
  -s string
    	Bus stop cache duration (default "168h")
  -x	Allow requests from other domains
```

## API

### `/`

Lists all available API routes.

Example:

```
$ curl https://mpolden.no/atb/ | jq .
{
  "urls": [
    "https://mpolden.no/atb/v2/departures"
  ]
}
```

### `/api/v2/departures`

List departures from the given bus stop, identified by a stop ID. Use
https://stoppested.entur.org to find stop IDs, for example `41613` (the number
part of `NSR:StopPlace:41613`) for Prinsens gate.

Departures traveling in any direction are included by default. Add the parameter
`direction=inbound` or `direction=outbound` to filter departures towards, or
away from, the city centre.

Note that the claimed direction is questionable in some cases so inspect the
responses to decide whether `inbound` or `outbound` makes sense for your use
case.

```
$ curl 'https://mpolden.no/atb/v2/departures/41613?direction=inbound' | jq .

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
