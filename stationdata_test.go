package dbapi

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestStationDataAPI_ByID(t *testing.T) {
	assert := assert.New(t)

	http.HandleFunc("/stada/v2/stations/1", func(writer http.ResponseWriter, request *http.Request) {
		filename := "testdata" + request.URL.Path + ".json"
		dat, _ := ioutil.ReadFile(filename)

		fmt.Fprint(writer, string(dat))
	})

	once.Do(startMockServer)

	APIURL = "http://" + serverAddr + "/"

	c := New("SomeFakeToken", Config{})
	s := c.StationDataAPI()

	stationResp, _ := s.ByID(1)
	station := stationResp.Result[0]

	assert.Equal(1, stationResp.Total)
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
