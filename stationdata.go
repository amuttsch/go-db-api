package dbapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/go-querystring/query"
)

const stadaAPIPath = "/stada/v2"

// MailingAddress holds the postal address of the station.
type MailingAddress struct {
	City        string `json:"city,omitempty"`
	Zipcode     string `json:"zipcode,omitempty"`
	Street      string `json:"street,omitempty"`
	HouseNumber string `json:"houseNumber,omitempty"`
}

// GeographicCoordinates holds the type of the coordinate and the latitude and longitude of the station.
type GeographicCoordinates struct {
	Type        string    `json:"type,omitempty"`
	Coordinates []float64 `json:"coordinates,omitempty"`
}

// EvaNumbers hold the eva (unique train station number), an identifier if this is the main one in case
// a station has multiple eva numbers and the coordinates of the station.
type EvaNumbers struct {
	Number                int                   `json:"number,omitempty"`
	IsMain                bool                  `json:"isMain,omitempty"`
	GeographicCoordinates GeographicCoordinates `json:"geographicCoordinates,omitempty"`
}

// Ril100Identifiers hold the ril (guideline, alphanumeric), an identifier if this is the main one in case
// a station has multiple eva numbers and the coordinates of the station.
type Ril100Identifiers struct {
	RilIdentifier         string                `json:"rilIdentifier,omitempty"`
	IsMain                bool                  `json:"isMain,omitempty"`
	HasSteamPermission    bool                  `json:"hasSteamPermission,omitempty"`
	GeographicCoordinates GeographicCoordinates `json:"geographicCoordinates,omitempty"`
}

// TimetableOffice holds information about the office that is responsible for the timetable at this station.
type TimetableOffice struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

// StationManagement holds information about the office that manages this station.
type StationManagement struct {
	Name   string `json:"name,omitempty"`
	Number int    `json:"number,omitempty"`
}

// OpeningTimes holds the opening and closing time.
type OpeningTimes struct {
	FromTime string `json:"fromTime,omitempty"`
	ToTime   string `json:"toTime,omitempty"`
}

// LocalServiceStaff holds times when local staff is available.
type LocalServiceStaff struct {
	Availability Availability `json:"availability,omitempty"`
}

// DBinformation holds opening times of the DB Information center at the station.
type DBinformation struct {
	Availability Availability `json:"availability,omitempty"`
}

// Availability holds the opening times for entities on various days.
type Availability struct {
	Monday    OpeningTimes `json:"monday,omitempty"`
	Tuesday   OpeningTimes `json:"tuesday,omitempty"`
	Wednesday OpeningTimes `json:"wednesday,omitempty"`
	Thursday  OpeningTimes `json:"thursday,omitempty"`
	Friday    OpeningTimes `json:"friday,omitempty"`
	Saturday  OpeningTimes `json:"saturday,omitempty"`
	Sunday    OpeningTimes `json:"sunday,omitempty"`
	Holiday   OpeningTimes `json:"holiday,omitempty"`
}

// Regionalbereich of the station.
// See: https://de.wikipedia.org/wiki/DB_Netz#Netzzugang
type Regionalbereich struct {
	Name      string `json:"name,omitempty"`
	ShortName string `json:"shortName,omitempty"`
	Number    int    `json:"number,omitempty"`
}

// SZentrale holds information about the 3S-Zentrale that is responsible for service, security and cleanliness for
// a station.
// See: https://www.bahnhof.de/bahnhof-de/ueberuns/3-s-konzept-519192
type SZentrale struct {
	Number            int    `json:"number,omitempty"`
	Name              string `json:"name,omitempty"`
	PublicPhoneNumber string `json:"publicPhoneNumber,omitempty"`
}

// Aufgabentraeger holds information about the entity that is responsible for local trains.
// See: https://www.dbregio.de/db_regio/view/wir/nahverkehr-deutschland.shtml
type Aufgabentraeger struct {
	Shortname string `json:"shortname,omitempty"`
	Name      string `json:"name,omitempty"`
}

// Station is a struct containing all fields that the StationData API might return.
type Station struct {
	Number                  int                 `json:"number,omitempty"`
	Name                    string              `json:"name,omitempty"`
	MailingAddress          MailingAddress      `json:"mailingAddress,omitempty"`
	Category                int                 `json:"category,omitempty"`
	PriceCategory           int                 `json:"priceCategory,omitempty"`
	FederalState            string              `json:"federalState,omitempty"`
	HasParking              bool                `json:"hasParking,omitempty"`
	HasBicycleParking       bool                `json:"hasBicycleParking,omitempty"`
	HasLocalPublicTransport bool                `json:"hasLocalPublicTransport,omitempty"`
	HasPublicFacilities     bool                `json:"hasPublicFacilities,omitempty"`
	HasLockerSystem         bool                `json:"hasLockerSystem,omitempty"`
	HasTaxiRank             bool                `json:"hasTaxiRank,omitempty"`
	HasTravelNecessities    bool                `json:"hasTravelNecessities,omitempty"`
	HasSteplessAccess       string              `json:"hasSteplessAccess,omitempty"`
	HasMobilityService      string              `json:"hasMobilityService,omitempty"`
	HasWiFi                 bool                `json:"hasWiFi,omitempty"`
	HasTravelCenter         bool                `json:"hasTravelCenter,omitempty"`
	HasRailwayMission       bool                `json:"hasRailwayMission,omitempty"`
	HasDBLounge             bool                `json:"hasDBLounge,omitempty"`
	HasLostAndFound         bool                `json:"hasLostAndFound,omitempty"`
	HasCarRental            bool                `json:"hasCarRental,omitempty"`
	EvaNumbers              []EvaNumbers        `json:"evaNumbers,omitempty"`
	Ril100Identifiers       []Ril100Identifiers `json:"ril100Identifiers,omitempty"`
	TimetableOffice         TimetableOffice     `json:"timetableOffice,omitempty"`
	StationManagement       StationManagement   `json:"stationManagement,omitempty"`
	LocalServiceStaff       LocalServiceStaff   `json:"localServiceStaff,omitempty"`
	DBinformation           DBinformation       `json:"DBinformation,omitempty"`
	Regionalbereich         Regionalbereich     `json:"regionalbereich,omitempty"`
	SZentrale               SZentrale           `json:"szentrale,omitempty"`
	Aufgabentraeger         Aufgabentraeger     `json:"aufgabentraeger,omitempty"`
}

// StationDataResponse holds meta information about the response and the actual station set.
type StationDataResponse struct {
	Offset int       `json:"offset,omitempty"`
	Total  int       `json:"total,omitempty"`
	Limit  int       `json:"limit,omitempty"`
	Result []Station `json:"result,omitempty"`
}

// StationDataRequest is used by ByFilter to query the station API. If it's not chaned,
// all stations are queried.
type StationDataRequest struct {
	Offset          int    `url:"offset,omitempty"`
	Limit           int    `url:"limit,omitempty"`
	Searchstring    string `url:"searchstring,omitempty"`
	Category        string `url:"category,omitempty"`
	Federalstate    string `url:"federalstate,omitempty"`
	Eva             int    `url:"eva,omitempty"`
	Ril             string `url:"ril,omitempty"`
	Logicaloperator string `url:"logicaloperator,omitempty"`
}

// StationDataRateErrorDetailsResponse contains information about the rate limit.
type StationDataRateErrorDetailsResponse struct {
	Code        int    `json:"code,omitempty"`
	Message     string `json:"message,omitempty"`
	Description string `json:"description,omitempty"`
}

// StationDataRateErrorResponse holds a field containing information about the rate limit. Occurs when the
// API was called too often.
type StationDataRateErrorResponse struct {
	Err StationDataRateErrorDetailsResponse `json:"error,omitempty"`
}

// StationDataErrorResponse holds error information returned by the DB API
type StationDataErrorResponse struct {
	ErrNo  int    `json:"errNo,omitempty"`
	ErrMsg string `json:"errMsg,omitempty"`
}

// StationDataAPI is a struct holding internal information about this API. Its methods can be used
// to query the API.
type StationDataAPI struct {
	client                *Client
	rateThrottleTicker    *time.Ticker
	firstRequestProcessed bool
}

func (e *StationDataRateErrorResponse) Error() string {
	return fmt.Sprintf("Error %d: %s - %s", e.Err.Code, e.Err.Message, e.Err.Description)
}

func (e *StationDataErrorResponse) Error() string {
	return fmt.Sprintf("Error %d: %s", e.ErrNo, e.ErrMsg)
}

// ByID returns station information for the given id or an error if the
// id is invalid, rate limiting or some other error occurred.
func (s *StationDataAPI) ByID(id int) (*StationDataResponse, error) {
	url := fmt.Sprintf("%s%s/stations/%d", APIURL, stadaAPIPath, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := s.sendRequest(req)
	if err != nil {
		return nil, err
	}

	return s.processResponse(resp)
}

// ByFilter returns a list of station information by the given filter or an error if the
// id is invalid, rate limiting or some other error occurred. If the StationDataRequest is
// not set, all stations are returned (max 10.000) - same as All().
func (s *StationDataAPI) ByFilter(stationRequest StationDataRequest) (*StationDataResponse, error) {
	q, err := query.Values(stationRequest)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s%s/stations?%s", APIURL, stadaAPIPath, q.Encode())

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := s.sendRequest(req)
	if err != nil {
		return nil, err
	}

	return s.processResponse(resp)
}

// All returns station information for all available stations. Same as calling
// ByFilter(StationDataRequest{}).
func (s *StationDataAPI) All() (*StationDataResponse, error) {
	return s.ByFilter(StationDataRequest{})
}

func (s *StationDataAPI) limitRate() {
	// Throttle API in case a tier was specified
	if s.rateThrottleTicker != nil && s.firstRequestProcessed {
		<-s.rateThrottleTicker.C
	}
	s.firstRequestProcessed = true
}

func (s *StationDataAPI) sendRequest(req *http.Request) (*http.Response, error) {
	if s.client.APIToken == "" {
		return nil, errors.New("no API token given")
	}

	s.limitRate()

	req.Header.Set("Authorization", "Bearer "+s.client.APIToken)
	return s.client.httpClient.Do(req)
}

func (s *StationDataAPI) processResponse(resp *http.Response) (sdr *StationDataResponse, err error) {
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	switch resp.StatusCode {
	case 200:
		stationDataResponse := StationDataResponse{}
		err := json.NewDecoder(resp.Body).Decode(&stationDataResponse)
		if err != nil {
			return nil, err
		}
		return &stationDataResponse, nil
	case 404, 500:
		stationDataErrorResponse := StationDataErrorResponse{}
		err := json.NewDecoder(resp.Body).Decode(&stationDataErrorResponse)
		if err != nil {
			return nil, err
		}
		return nil, &stationDataErrorResponse
	case 429:
		stationDataRateErrorResponse := StationDataRateErrorResponse{}
		err := json.NewDecoder(resp.Body).Decode(&stationDataRateErrorResponse)
		if err != nil {
			return nil, err
		}
		return nil, &stationDataRateErrorResponse
	default:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(string(body))
	}
}
