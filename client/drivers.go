package client

import (
	"log"
    "encoding/xml"
)

type DriversMRData struct {
	XMLName        xml.Name             `xml:"MRData"`
	StandingsTable DriverStandingsTable `xml:"StandingsTable"`
}

type DriverStandingsTable struct {
	XMLName       xml.Name            `xml:"StandingsTable"`
	StandingsList DriverStandingsList `xml:"StandingsList"`
}

type DriverStandingsList struct {
	XMLName        xml.Name          `xml:"StandingsList"`
	DriverStanding []DriverStanding `xml:"DriverStanding"`
}

type DriverStanding struct {
	XMLName      xml.Name          `xml:"DriverStanding"`
	Position     int               `xml:"position,attr"`
	PositionText string            `xml:"positionText,attr"`
	Points       string               `xml:"points,attr"`
	Wins         int               `xml:"wins,attr"`
	Driver       Driver            `xml:"Driver"`
	Constructor  DriverConstructor `xml:"Constructor"`
}

type Driver struct {
	XMLName       xml.Name `xml:"Driver"`
	DriverID      string   `xml:"driverId,attr"`
	Code          string   `xml:"code,attr"`
	URL           string   `xml:"url,attr"`
	PermanentNum  int      `xml:"PermanentNumber"`
	GivenName     string   `xml:"GivenName"`
	FamilyName    string   `xml:"FamilyName"`
	DateOfBirth   string   `xml:"DateOfBirth"`
	Nationality   string   `xml:"Nationality"`
}

type DriverConstructor struct {
	XMLName       xml.Name `xml:"Constructor"`
	ConstructorID string   `xml:"constructorId,attr"`
	URL           string   `xml:"url,attr"`
	Name          string   `xml:"Name"`
	Nationality   string   `xml:"Nationality"`
}

func DriversParseData(requestBody []byte) string {
    var driversData DriversMRData
	response := "```\n"
    err := xml.Unmarshal(requestBody, &driversData)
    if err != nil {
        log.Println(err)
    }

	for _, elem := range driversData.StandingsTable.StandingsList.DriverStanding {
		response = response + elem.Driver.Code + " " + elem.Points + "\n"
	}
	response = response + "```"
	return response
}
