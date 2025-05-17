package client

import (
	"encoding/json"
	"log"
	"strings"
)

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
