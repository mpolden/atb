# atbapi

[![Build Status](https://travis-ci.org/martinp/atbapi.png)](https://travis-ci.org/martinp/atbapi)

JSON API for bus data in Trondheim, Norway. This API proxies requests to the AtB
public API and converts the responses into a sane format.

You can request access to the SOAP API provided by AtB
[here](https://www.atb.no/aapne-data/category419.html) (Norwegian).

## Usage

```
$ atbapi -h
Usage:
  atbapi [OPTIONS]

Application Options:
  -l, --listen=        Listen address (:8080)
  -c, --config=FILE    Path to config file (config.json)

Help Options:
  -h, --help           Show this help message
```

## Example config

```
{
  "Username": "username",
  "Password": "password"
}
```





