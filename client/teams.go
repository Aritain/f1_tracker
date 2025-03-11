package client

import (
	"encoding/json"
	"log"
	"strings"
)

type TeamsResponse struct {
	MRData TeamsMRData `json:"MRData"`
}

type TeamsMRData struct {
	XMLNS          string             `json:"xmlns"`
	Series         string             `json:"series"`
	URL            string             `json:"url"`
	Limit          string             `json:"limit"`
	Offset         string             `json:"offset"`
	Total          string             `json:"total"`
	StandingsTable TeamStandingsTable `json:"StandingsTable"`
}

type TeamStandingsTable struct {
	Season         string               `json:"season"`
	Round          string               `json:"round"`
	StandingsLists []TeamStandingsLists `json:"StandingsLists"`
}

type TeamStandingsLists struct {
	Season               string                 `json:"season"`
	Round                string                 `json:"round"`
	ConstructorStandings []ConstructorStandings `json:"ConstructorStandings"`
}

type ConstructorStandings struct {
	Position     string          `json:"position"`
	PositionText string          `json:"positionText"`
	Points       string          `json:"points"`
	Wins         string          `json:"wins"`
	Constructor  TeamConstructor `json:"Constructor"`
}

type TeamConstructor struct {
	ConstructorID string `json:"constructorId"`
	URL           string `json:"url"`
	Name          string `json:"name"`
	Nationality   string `json:"nationality"`
}

func TeamsParseData(requestBody []byte) string {
	var teamsData TeamsResponse
	var longestTeam int

	response := "```\n"
	err := json.Unmarshal(requestBody, &teamsData)
	if err != nil {
		log.Println(err)
	}
	if len(teamsData.MRData.StandingsTable.StandingsLists) == 0 {
		response = "Team information is not available yet. Try again later 🙃"
		return response
	}
	/*
	   Make team names a bit nicer by removing unnececary keywords, e.g. "Alpine F1 Team" -> "Alpine"
	   Magic [0] is current year, the API returns a list of StandingsLists but it always contains only 1 element
	*/
	for index, elem := range teamsData.MRData.StandingsTable.StandingsLists[0].ConstructorStandings {
		teamsData.MRData.StandingsTable.StandingsLists[0].ConstructorStandings[index].Constructor.Name =
			strings.ReplaceAll(teamsData.MRData.StandingsTable.StandingsLists[0].ConstructorStandings[index].Constructor.Name, "Team", "")
		teamsData.MRData.StandingsTable.StandingsLists[0].ConstructorStandings[index].Constructor.Name =
			strings.ReplaceAll(teamsData.MRData.StandingsTable.StandingsLists[0].ConstructorStandings[index].Constructor.Name, "F1", "")
		teamsData.MRData.StandingsTable.StandingsLists[0].ConstructorStandings[index].Constructor.Name =
			strings.ReplaceAll(teamsData.MRData.StandingsTable.StandingsLists[0].ConstructorStandings[index].Constructor.Name, "Team", "")
		if longestTeam < len(elem.Constructor.Name) {
			longestTeam = len(elem.Constructor.Name)
		}
	}

	for _, elem := range teamsData.MRData.StandingsTable.StandingsLists[0].ConstructorStandings {
		spaces := strings.Repeat(" ", (longestTeam - len(elem.Constructor.Name)))
		response = response + elem.Constructor.Name + spaces + elem.Points + "\n"
	}

	response = response + "```"
	return response
}
