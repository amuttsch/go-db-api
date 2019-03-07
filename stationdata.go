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

type MailingAddress struct {
	City        string `json:"city,omitempty"`
	Zipcode     string `json:"zipcode,omitempty"`
	Street      string `json:"street,omitempty"`
	HouseNumber string `json:"houseNumber,omitempty"`
}

type GeographicCoordinates struct {
	Type        string    `json:"type,omitempty"`
	Coordinates []float64 `json:"coordinates,omitempty"`
}

type EvaNumbers struct {
	Number                int                   `json:"number,omitempty"`
	IsMain                bool                  `json:"isMain,omitempty"`
	GeographicCoordinates GeographicCoordinates `json:"geographicCoordinates,omitempty"`
}

type Ril100Identifiers struct {
	RilIdentifier         string                `json:"rilIdentifier,omitempty"`
	IsMain                bool                  `json:"isMain,omitempty"`
	HasSteamPermission    bool                  `json:"hasSteamPermission,omitempty"`
	GeographicCoordinates GeographicCoordinates `json:"geographicCoordinates,omitempty"`
}

type TimetableOffice struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

type StationManagement struct {
	Name   string `json:"name,omitempty"`
	Number int    `json:"number,omitempty"`
}

type OpeningTimes struct {
	FromTime string `json:"fromTime,omitempty"`
	ToTime   string `json:"toTime,omitempty"`
}

type LocalServiceStaff struct {
	Availability Availability `json:"availability,omitempty"`
}

type DBinformation struct {
	Availability Availability `json:"availability,omitempty"`
}

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

type Regionalbereich struct {
	Name      string `json:"name,omitempty"`
	ShortName string `json:"shortName,omitempty"`
	Number    int    `json:"number,omitempty"`
}

type SZentrale struct {
	Number            int    `json:"number,omitempty"`
	Name              string `json:"name,omitempty"`
	PublicPhoneNumber string `json:"publicPhoneNumber,omitempty"`
}

type Aufgabentraeger struct {
	Shortname string `json:"shortname,omitempty"`
	Name      string `json:"name,omitempty"`
}

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

type StationDataResponse struct {
	Offset int       `json:"offset,omitempty"`
	Total  int       `json:"total,omitempty"`
	Limit  int       `json:"limit,omitempty"`
	Result []Station `json:"result,omitempty"`
}

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

type StationDataRateErrorDetailsResponse struct {
	Code        int    `json:"code,omitempty"`
	Message     string `json:"message,omitempty"`
	Description string `json:"description,omitempty"`
}

type StationDataRateErrorResponse struct {
	Err StationDataRateErrorDetailsResponse `json:"error,omitempty"`
}

type StationDataErrorResponse struct {
	ErrNo  int    `json:"errNo,omitempty"`
	ErrMsg string `json:"errMsg,omitempty"`
}

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
func (s *StationDataAPI) ByFilter(stationRequest *StationDataRequest) (*StationDataResponse, error) {
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
	return s.ByFilter(&StationDataRequest{})
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
