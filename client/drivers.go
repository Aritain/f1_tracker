package client

import (
    "log"
    "encoding/json"
)

type DriversResponse struct {
	MRData DriverMRData `json:"MRData"`
}

type DriverMRData struct {
	XMLNS          string         `json:"xmlns"`
	Series         string         `json:"series"`
	URL            string         `json:"url"`
	Limit          string         `json:"limit"`
	Offset         string         `json:"offset"`
	Total          string         `json:"total"`
	StandingsTable DriverStandingsTable `json:"StandingsTable"`
}

type DriverStandingsTable struct {
	Season         string          `json:"season"`
	Round          string          `json:"round"`
	StandingsList []DriverStandingsList `json:"StandingsLists"`
}

type DriverStandingsList struct {
	Season          string            `json:"season"`
	Round           string            `json:"round"`
	DriverStanding []DriverStanding `json:"DriverStandings"`
}

type DriverStanding struct {
	Position      string       `json:"position"`
	PositionText  string       `json:"positionText"`
	Points        string       `json:"points"`
	Wins          string       `json:"wins"`
	Driver        Driver       `json:"Driver"`
	Constructors  []DriverConstructor `json:"Constructors"`
}

type Driver struct {
	DriverID         string `json:"driverId"`
	PermanentNumber  string `json:"permanentNumber"`
	Code            string `json:"code"`
	URL             string `json:"url"`
	GivenName       string `json:"givenName"`
	FamilyName      string `json:"familyName"`
	DateOfBirth     string `json:"dateOfBirth"`
	Nationality     string `json:"nationality"`
}

type DriverConstructor struct {
	ConstructorID string `json:"constructorId"`
	URL          string `json:"url"`
	Name         string `json:"name"`
	Nationality  string `json:"nationality"`
}

func DriversParseData(requestBody []byte) string {
    var driversData DriversResponse
    response := "```\n"
    err := json.Unmarshal(requestBody, &driversData)
    if err != nil {
        log.Println(err)
    }
    if len(driversData.MRData.StandingsTable.StandingsList) == 0 {
        response = "Driver information is not available yet. Try again later 🙃"
        return response
    }
    //  Magic [0] is current year, the API returns a list of StandingsLists but it always contains only 1 element
    for _, elem := range driversData.MRData.StandingsTable.StandingsList[0].DriverStanding {
        response = response + elem.Driver.Code + " " + elem.Points + "\n"
    }
    response = response + "```"
    return response
}
