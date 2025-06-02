package client

import (
	"log"
	"strconv"
	"time"
)

const DATE_FORMAT = "2006-01-02 15:04:05Z"

func RaceUpdater() {
	params := map[string]string{}
	headers := map[string]string{}
	var APIResponse RacesResponse
	var err error
	var requestUrl string
	var raceDate string
	var storedRace ParsedRace
	var candidateRace ParsedRace
	for {
		year, _, _ := time.Now().Date()
		now := time.Now()
		// For tests
		//t := "2025-05-01T13:00:00Z"
		//now, _ := time.Parse(time.RFC3339, t)
		requestUrl = BASE_URL + strconv.Itoa(year) + RACES_URL

		APIResponse, err = GetRequest[RacesResponse](
			requestUrl,
			"json",
			params,
			headers,
		)
		if err != nil {
			log.Println(err)
			time.Sleep(1 * time.Hour)
			continue
		}
		for _, race := range APIResponse.MRData.RaceTable.Race {
			raceDate = CombineDate(race.Date, race.Time)
			dateParsed, _ := time.Parse(DATE_FORMAT, raceDate)
			if dateParsed.After(now) {
				candidateRace = ParseRase(race)
				break
			}
		}
		storedRace = LoadRace()
		if storedRace.RaceDate != candidateRace.RaceDate {
			WriteRace(candidateRace)
		}
		time.Sleep(1 * time.Hour)

	}
}

func RaceChecker() {
	var msg string
	for {
		race := LoadRace()
		// Skip further checks if particular event was already notified of
		if race.SprintTrigger && race.RaceTrigger {
			time.Sleep(30 * time.Minute) // Wait 30 minutes
			continue
		}
		// Add two hours to time.now to notify about the event 2 hours in advance
		offsetNow := time.Now().Add(2 * time.Hour)
		// For tests
		//t := "2025-06-15T18:00:00Z"
		//offsetNow, _ := time.Parse(time.RFC3339, t)
		sprintCheck := isSameDateTime(race.SprintDate.UTC(), offsetNow.UTC())
		raceCheck := isSameDateTime(race.RaceDate.UTC(), offsetNow.UTC())

		if sprintCheck && !race.SprintTrigger {
			race.SprintTrigger = true
			WriteRace(race)
			msg = PrepRaceMessage(race, "sprint")
			SendTGMessage(msg)
			continue
		}
		if raceCheck && !race.RaceTrigger {
			race.RaceTrigger = true
			WriteRace(race)
			msg = PrepRaceMessage(race, "race")
			SendTGMessage(msg)
			continue
		}

		time.Sleep(15 * time.Second) // Wait 15 seconds
	}
}
