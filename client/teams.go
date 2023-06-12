package client

import (
    "log"
    "strings"
    "encoding/xml"
)


type TeamsMRData struct {
    XMLName        xml.Name           `xml:"MRData"`
    StandingsTable TeamStandingsTable `xml:"StandingsTable"`
}

type TeamStandingsTable struct {
    XMLName       xml.Name          `xml:"StandingsTable"`
    StandingsList TeamStandingsList `xml:"StandingsList"`
}

type TeamStandingsList struct {
    XMLName             xml.Name              `xml:"StandingsList"`
    ConstructorStanding []ConstructorStanding `xml:"ConstructorStanding"`
}

type ConstructorStanding struct {
    XMLName      xml.Name        `xml:"ConstructorStanding"`
    Position     int             `xml:"position,attr"`
    PositionText string          `xml:"positionText,attr"`
    Points       string          `xml:"points,attr"`
    Wins         int             `xml:"wins,attr"`
    Constructor  TeamConstructor `xml:"Constructor"`
}

type TeamConstructor struct {
    XMLName       xml.Name `xml:"Constructor"`
    ConstructorID string   `xml:"constructorId,attr"`
    URL           string   `xml:"url,attr"`
    Name          string   `xml:"Name"`
    Nationality   string   `xml:"Nationality"`
}

func TeamsParseData(requestBody []byte) string {
    var teamsData TeamsMRData
    var longestTeam int

    response := "```\n"
    err := xml.Unmarshal(requestBody, &teamsData)
    if err != nil {
        log.Println(err)
    }

    // Make team names a bit nicer by removing unnececary keywords
    for index, elem := range teamsData.StandingsTable.StandingsList.ConstructorStanding {
        teamsData.StandingsTable.StandingsList.ConstructorStanding[index].Constructor.Name =
            strings.ReplaceAll(teamsData.StandingsTable.StandingsList.ConstructorStanding[index].Constructor.Name, "Team", "")
        teamsData.StandingsTable.StandingsList.ConstructorStanding[index].Constructor.Name =
            strings.ReplaceAll(teamsData.StandingsTable.StandingsList.ConstructorStanding[index].Constructor.Name, "F1", "")
        teamsData.StandingsTable.StandingsList.ConstructorStanding[index].Constructor.Name =
            strings.ReplaceAll(teamsData.StandingsTable.StandingsList.ConstructorStanding[index].Constructor.Name, "Team", "")
        if longestTeam < len(elem.Constructor.Name) {
            longestTeam = len(elem.Constructor.Name)
        }
    }

    for _, elem := range teamsData.StandingsTable.StandingsList.ConstructorStanding {
        spaces := strings.Repeat(" ", (longestTeam - len(elem.Constructor.Name)))
        response = response + elem.Constructor.Name + spaces + elem.Points + "\n"
    }

    response = response + "```"
    return response
}
