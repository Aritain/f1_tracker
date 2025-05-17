package client

import (
	"encoding/json"
	"log"
)

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
