// Package dbapi provides access to the open data APIs from Deutsche Bahn.
//
// See https://developer.deutschebahn.com/
package dbapi

import (
	"net/http"
	"sync"
	"time"
)

// APIURL is made public for testing reasons
var APIURL = "https://api.deutschebahn.com"

// Client enables access to the open data APIs provided by Deutsche Bahn. An access token
// is required to query the APIs and can be obtained for free at https://developer.deutschebahn.com/store/site/pages/sign-up.jag
//
// This struct enables access to all implemented APIs by providing methods to get the API structs,
// e.g. Client.StationDataAPI().
//
// See https://developer.deutschebahn.com/store/apis/list for a complete list of
// available APIs.
type Client struct {
	APIToken   string
	httpClient *http.Client
	apiConfig  Config

	stationDataAPI            *StationDataAPI
	stationDataAPIInitialized sync.Once
}

// StationDataConfig provides configuration options for the StationData API. Set rateLimitPerMinute to
// zero if you want to disable rate limiting done in the library.
type StationDataConfig struct {
	rateLimitPerMinute int
}

// Config provides configuration for all implemented APIs.
type Config struct {
	StationDataConfig StationDataConfig
}

// New creates a new Client and needs an API token. It provides access to all implemented
// APIs and handles rate limiting (if set in Config).
func New(token string, apiConfig Config) *Client {
	return &Client{
		APIToken: token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiConfig: apiConfig,
	}
}

// StationDataAPI provides access to the StationData v2 API located at https://developer.deutschebahn.com/store/apis/info?name=StaDa-Station_Data&version=v2&provider=DBOpenData
// It is possible to query Stations and 3S-central points either by filter or by id.
func (client *Client) StationDataAPI() *StationDataAPI {
	client.stationDataAPIInitialized.Do(func() {
		var ticker *time.Ticker
		rateLimitPerMinute := float64(client.apiConfig.StationDataConfig.rateLimitPerMinute)
		if rateLimitPerMinute > 0 {
			rate := time.Duration(60.0/rateLimitPerMinute*1000) * time.Millisecond
			ticker = time.NewTicker(rate)
		}
		client.stationDataAPI = &StationDataAPI{
			client:             client,
			rateThrottleTicker: ticker,
		}
	})

	return client.stationDataAPI
}
