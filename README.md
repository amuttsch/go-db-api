# DB Open Data in Go 
[![GoDoc](https://godoc.org/github.com/amuttsch/go-db-api?status.svg)](https://godoc.org/github.com/amuttsch/go-db-api) [![Build Status](https://travis-ci.org/amuttsch/go-db-api.svg?branch=master)](https://travis-ci.org/amuttsch/go-db-api) [![Go Report Card](https://goreportcard.com/badge/github.com/amuttsch/go-db-api)](https://goreportcard.com/report/github.com/amuttsch/go-db-api) [![License](https://img.shields.io/badge/License-BSD%202--Clause-orange.svg)](https://opensource.org/licenses/BSD-2-Clause)

This client library enables you to easily access supported APIs available from [Deutsche Bahn Open API](https://developer.deutschebahn.com/store/site/pages/home.jag) from any Go application.

Currently supported:

| API | Status
|-----|------------
| [Station Data v2](https://developer.deutschebahn.com/store/apis/info?name=StaDa-Station_Data&version=v2&provider=DBOpenData)    | Partial (stations) |

## Installation

    $ go get -u github.com/amuttsch/go-db-api
    
Or via Go modules import

    import (
       github.com/amuttsch/go-db-api
    )

In order to use the API you have to [sign up](https://developer.deutschebahn.com/store/site/pages/sign-up.jag) for a free API token and subscribe to the APIs (free).

## Example

Get all stations in the federal state of Hessen:

    api := New("your token", APIConfig{})
    stationDataAPI := api.StationDataAPI()

    stationResponse, err := stationDataAPI.ByFilter(&StationDataRequest{
        Federalstate: "hessen",
    })
    
    for station := stationResponse.Result {
        fmt.Println(station.Name)
    }


## Rate limiting

Most APIs from Deutsche Bahn are rate limited. When you subscribe to an API you have to choose a tier which sets the amount of requests you can make on this API. `go-db-api` has a built in rate limiting which blocks until the next request can be made if you configure it in the `APIConfig`. In the case of a limit of 10 requests per minute, each 6 seconds a request is allowed to process.

If you want to implement your own rate limiting, set the `rateLimitPerMinute` to zero (default).
