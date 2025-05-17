package client

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type RaceData struct {
	Location   string
	SprintInfo string
	Sprint     string
	SprintES   string
	RaceInfo   string
	Race       string
	RaceES     string
}

const UPDATE_INTERVAL = 1
const FILEPATH = "/data/race.json"

// Extremely cursed code, need to rework & maybe use pointers
// Ideally store next race data locally and fetch that saved data upon user request
// And when processing the announcements

func AssetWatcher(bot *tgbotapi.BotAPI) {
	notificationID, status := os.LookupEnv("NOTIFICATION_ID")
	if !status {
		log.Printf("NOTIFICATION_ID env is missing.")
		os.Exit(1)
	}
	notificationIDNum, _ := strconv.ParseInt(notificationID, 10, 64)

	for {
		var raceParsed RaceData
		raceData := FetchData("next_race")
		if (raceData == "No more races this year 😭") || (strings.Contains(raceData, "To be announced")) {
			time.Sleep(1 * time.Hour)
			continue
		}
		raceData = strings.Replace(raceData, "\n🇬🇧", " ", -1)
		raceData = strings.Replace(raceData, "\n🇪🇸", " - 🇪🇸", -1)
		raceOutput := strings.Split(raceData, "\n")
		for index, elem := range raceOutput {
			switch {
			case elem == "":
				continue
			case index == 0:
				continue
			case index == 1:
				raceParsed.Location = elem
			case strings.Contains(elem, "Shootout"):
				continue
			case strings.Contains(elem, "Qualifying"):
				continue
			case strings.Contains(elem, "🏎 Sprint"):
				raceParsed.SprintInfo, raceParsed.Sprint, raceParsed.SprintES = CutData(elem, " - ")
			case strings.Contains(elem, "Race"):
				raceParsed.RaceInfo, raceParsed.Race, raceParsed.RaceES = CutData(elem, " - ")
			}
		}

		if raceParsed.Sprint != "" {
			if ProcessMsg(bot, raceParsed.Location, raceParsed.Sprint, raceParsed.SprintES, raceParsed.SprintInfo, notificationIDNum) {
				time.Sleep(1 * time.Hour)
				continue
			}
		}

		ProcessMsg(bot, raceParsed.Location, raceParsed.Race, raceParsed.RaceES, raceParsed.RaceInfo, notificationIDNum)
		time.Sleep(1 * time.Hour)
	}
}

func ProcessMsg(bot *tgbotapi.BotAPI, eventLocation string, eventDateRaw string, esTime string, eventType string, notificationIDNum int64) bool {
	currentDate := time.Now()
	year, _, _ := time.Now().Date()
	// Uncomment for testing
	//currentDate, _ := time.Parse("2006-01-02 15:04:05", "2024-11-30 15:30:00")
	//currentDate, _ := time.Parse("2006-01-02 15:04:05", "2024-11-29 13:30:00")
	//year := 2024

	eventDate, _ := time.Parse("January 02 15:04", eventDateRaw)
	eventDate = eventDate.AddDate(year-eventDate.Year(), 0, 0)

	if CheckDates(currentDate, eventDate) {
		_, _, gbTime := CutData(eventDateRaw, " ")
		msgText := eventType + " tomorrow in\n" + eventLocation + "\n" + "🇬🇧" + gbTime + "\n" + esTime
		msg := tgbotapi.NewMessage(notificationIDNum, msgText)
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
		return true
	}
	return false
}

func CutData(inputData string, pattern string) (string, string, string) {
	var tempSlice []string
	var a string
	var b string
	var c string

	tempSlice = strings.Split(inputData, pattern)
	a = tempSlice[0]
	b = tempSlice[1]
	c = tempSlice[2]
	return a, b, c
}

func CheckDates(currentDate time.Time, futureDate time.Time) bool {
	diff := futureDate.Sub(currentDate)
	minutes := diff.Minutes()
	if (minutes > 1440) && (minutes < 1500) {
		return true
	}
	return false
}

func RaceUpdater() {
	params := map[string]string{}
	headers := map[string]string{}
	var APIResponse RacesResponse
	var err error
	var requestUrl string
	var raceDate string
	var dateParsed time.Time
	var storedRace StoredRace
	var candidateRace StoredRace
	for {
		year, _, _ := time.Now().Date()
		requestUrl = BASE_URL + strconv.Itoa(year) + RACES_URL

		APIResponse, err = GetRequest[RacesResponse](
			requestUrl,
			"xml",
			params,
			headers,
		)
		if err != nil {
			log.Println(err)
			time.Sleep(UPDATE_INTERVAL * time.Hour)
			continue
		}
		for _, race := range APIResponse.MRData.RaceTable.Race {
			raceDate = race.Date + race.Time
			dateParsed, _ := time.Parse(time.RFC3339, raceDate)
			if dateParsed.After(time.Now()) {
				break
			}
		}
		storedRace = LoadRace()
		candidateRace = StoredRace{Date: dateParsed}
		if storedRace != candidateRace {
			WriteRace(candidateRace)
		}
		time.Sleep(UPDATE_INTERVAL * time.Hour)

	}
}

func RaceChecker() {
	for {
		time.Sleep(1 * time.Minute)
	}
}

func WriteRace(candidateRace StoredRace) {
	file, _ := os.Create(FILEPATH)
	defer file.Close()
	json.NewEncoder(file).Encode(candidateRace)
	log.Printf("New race written - %v.", candidateRace)
}

func LoadRace() (storedRace StoredRace) {
	data, err := os.ReadFile(FILEPATH)
	if err != nil {
		return
	}
	_ = json.Unmarshal(data, &storedRace)
	return
}
