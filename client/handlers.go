package client

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Fromat: BASE_URL + YEAR + SPECIFIC URL
const BASE_URL = "https://api.jolpi.ca/ergast/f1/"
const DRIVERS_URL = "/driverstandings/"
const RACES_URL = "/races/"
const TEAMS_URL = "/constructorstandings/"

func FetchData(userRequest string) string {
	var response string
	var requestUrl string
	year, _, _ := time.Now().Date()

	// Uncomment for testing
	//year := 2024
	requestUrl = BASE_URL + strconv.Itoa(year)
	switch {
	case userRequest == "teams":
		requestUrl = requestUrl + TEAMS_URL
	case userRequest == "drivers":
		requestUrl = requestUrl + DRIVERS_URL
	case userRequest == "next_race":
		requestUrl = requestUrl + RACES_URL
	default:
		response = "Command not recognised, try these:\n" +
			"/drivers - show drivers leaderboard\n" +
			"/teams - show teams leaderboard\n" +
			"/next\\_race - show next race information\n"
		return response
	}
	resp, err := http.Get(requestUrl)
	if err != nil {
		log.Println(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	switch {
	case userRequest == "teams":
		response = TeamsParseData(body)
	case userRequest == "drivers":
		response = DriversParseData(body)
	case userRequest == "next_race":
		response = RacesParseData(body)
	}

	return response
}
