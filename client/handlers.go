package client

import (
    "log"
    "io/ioutil"
    "net/http"
)

const TEAMS_URL = "http://ergast.com/api/f1/2024/constructorStandings"
const DRIVERS_URL = "http://ergast.com/api/f1/2024/driverStandings"
const RACES_URL ="http://ergast.com/api/f1/2024"

func FetchData (userRequest string) string {
    var response string
    var requestUrl string
    switch {
        case userRequest == "teams":
            requestUrl = TEAMS_URL
        case userRequest == "drivers":
            requestUrl = DRIVERS_URL
        case userRequest == "next_race":
            requestUrl = RACES_URL
        default:
            response = "Command not recognised, try these:\n" +
                "/drivers - show drivers leaderboard\n" +
                "/teams - show teams leaderboard\n" +
                "/next\\_race - show next race information\n"
            return response
    }

    resp, err := http.Get(requestUrl)
    if err != nil {
        log.Println("error")
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
