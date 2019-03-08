package dbapi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	http.HandleFunc("/stada/v2/", func(writer http.ResponseWriter, request *http.Request) {
		filename := "testdata" + request.URL.String() + ".json"

		fmt.Println("Sending data from file ", filename)

		dat, _ := ioutil.ReadFile(filename)

		fmt.Fprint(writer, string(dat))
	})
}

func TestStationDataAPI_StationByID(t *testing.T) {
	assert := assert.New(t)

	once.Do(startMockServer)

	APIURL = "http://" + serverAddr

	c := New("SomeFakeToken", Config{})
	s := c.StationDataAPI()

	stationResp, _ := s.StationByID(1)

	assert.NotNil(stationResp)

	assert.Equal(1, stationResp.Total)
	station := stationResp.Result[0]

	assert.Equal("Aachen Hbf", station.Name)
	assert.Equal(1, station.Number)
	assert.Equal(false, station.HasDBLounge)
	assert.Equal(true, station.HasParking)
	assert.Equal(true, station.HasBicycleParking)
	assert.Equal(true, station.HasWiFi)
	assert.Equal("Ja, um Voranmeldung unter 01806 512 512 wird gebeten", station.HasMobilityService)
	assert.Equal("Aachen", station.MailingAddress.City)
	assert.Equal("RB West", station.Regionalbereich.Name)
	assert.Equal("NVR", station.Aufgabentraeger.Shortname)
	assert.Equal("06:00", station.LocalServiceStaff.Availability.Monday.FromTime)
	assert.Equal("22:30", station.LocalServiceStaff.Availability.Monday.ToTime)
	assert.Equal("Bahnhofsmanagement Köln", station.TimetableOffice.Name)
	assert.Equal("Duisburg Hbf", station.SZentrale.Name)
	assert.Equal(45, station.StationManagement.Number)
	assert.Equal(8000001, station.EvaNumbers[0].Number)
	assert.Equal("Point", station.EvaNumbers[0].GeographicCoordinates.Type)
	assert.Equal(6.091499, station.EvaNumbers[0].GeographicCoordinates.Coordinates[0])
	assert.Equal(50.7678, station.EvaNumbers[0].GeographicCoordinates.Coordinates[1])
	assert.Equal("KA", station.Ril100Identifiers[0].RilIdentifier)
	assert.Equal(true, station.Ril100Identifiers[0].IsMain)
	assert.Equal(true, station.Ril100Identifiers[0].HasSteamPermission)
	assert.Equal("Point", station.Ril100Identifiers[0].GeographicCoordinates.Type)
	assert.Equal(6.091201396, station.Ril100Identifiers[0].GeographicCoordinates.Coordinates[0])
	assert.Equal(50.767558188, station.Ril100Identifiers[0].GeographicCoordinates.Coordinates[1])
}

func TestStationDataAPI_StationByFilter(t *testing.T) {
	assert := assert.New(t)

	once.Do(startMockServer)

	APIURL = "http://" + serverAddr

	c := New("SomeFakeToken", Config{})
	s := c.StationDataAPI()

	stationResp, _ := s.StationByFilter(StationDataStationRequest{
		Federalstate: "hessen",
	})
	assert.NotNil(stationResp)

	assert.Equal(429, stationResp.Total)
}

func TestStationDataAPI_SZentralenByID(t *testing.T) {
	assert := assert.New(t)

	once.Do(startMockServer)

	APIURL = "http://" + serverAddr

	c := New("SomeFakeToken", Config{})
	s := c.StationDataAPI()

	szResp, _ := s.SZentralenByID(15)

	assert.NotNil(szResp)

	assert.Equal(1, szResp.Total)
	sz := szResp.Result[0]

	assert.Equal("Duisburg Hbf", sz.Name)
}

func TestStationDataAPI_SZentralenAll(t *testing.T) {
	assert := assert.New(t)

	once.Do(startMockServer)

	APIURL = "http://" + serverAddr

	c := New("SomeFakeToken", Config{})
	s := c.StationDataAPI()

	szentralenResp, _ := s.SZentralenAll()
	assert.NotNil(szentralenResp)

	assert.Equal(30, szentralenResp.Total)
}

func TestRateLimiter(t *testing.T) {
	assert := assert.New(t)

	once.Do(startMockServer)

	APIURL = "http://" + serverAddr + "/"

	c := New("SomeFakeToken", Config{
		StationDataConfig: StationDataConfig{rateLimitPerMinute: 20}, // Should sleep for ~3 seconds
	})
	s := c.StationDataAPI()

	timeStartFirst := time.Now()
	s.StationByID(1)
	timeStartSecond := time.Now()
	s.StationByID(1)
	timeEnd := time.Now()

	durationFirstCall := timeStartSecond.Sub(timeStartFirst)
	durationSecondCall := timeEnd.Sub(timeStartSecond)

	assert.True(durationFirstCall.Seconds() < 1)
	assert.True(durationSecondCall.Seconds() > 3)
}
