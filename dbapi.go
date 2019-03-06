// Package dbapi provides access to the open data APIs from Deutsche Bahn.
//
// See https://developer.deutschebahn.com/
package dbapi

import (
	"net/http"
	"sync"
	"time"
)

// Make API url public for testing reasons
var APIURL = "https://api.deutschebahn.com/"

// APIClient enables access to the open data APIs provided by Deutsche Bahn. An access token
// is required to query the APIs and can be obtained for free at https://developer.deutschebahn.com/store/site/pages/sign-up.jag
//
// This struct enables access to all implemented APIs by providing methods to get the API structs,
// e.g. APIClient.StationDataAPI().
//
// See https://developer.deutschebahn.com/store/apis/list for a complete list of
// available APIs.
type APIClient struct {
	APIToken   string
	httpClient *http.Client
	apiConfig  APIConfig

	stationDataAPI            *StationDataAPI
	stationDataAPIInitialized sync.Once
}

// StationDataAPIConfig provides configuration options for the StationData API. Set rateLimitPerMinute to
// zero if you want to disable rate limiting done in the library.
type StationDataAPIConfig struct {
	rateLimitPerMinute int
}

// APIConfig provides configuration for all implemented APIs.
type APIConfig struct {
	StationDataAPIConfig StationDataAPIConfig
}

// New creates a new APIClient and needs an API token. It provides access to all implemented
// APIs and handles rate limiting (if set in APIConfig).
func New(token string, apiConfig APIConfig) *APIClient {
	return &APIClient{
		APIToken: token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiConfig: apiConfig,
	}
}

// StationDataAPI provides access to the StationData v2 API located at https://developer.deutschebahn.com/store/apis/info?name=StaDa-Station_Data&version=v2&provider=DBOpenData
// It is possible to query Stations and 3S-central points either by filter or by id.
func (client *APIClient) StationDataAPI() *StationDataAPI {
	client.stationDataAPIInitialized.Do(func() {
		var ticker *time.Ticker
		rateLimitPerMinute := float64(client.apiConfig.StationDataAPIConfig.rateLimitPerMinute)
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
