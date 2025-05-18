package client

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	t "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const FILEPATH = "/data/race.json"

func GetRequest[T any](url string, mode string, params map[string]string, headers map[string]string) (T, error) {
	var results T

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return results, err
	}

	reqParams := req.URL.Query()
	for k, v := range params {
		reqParams.Add(k, v)
	}
	req.URL.RawQuery = reqParams.Encode()

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return results, err
	}
	if mode == "xml" {
		err = xml.NewDecoder(resp.Body).Decode(&results)
	} else if mode == "json" {
		err = json.NewDecoder(resp.Body).Decode(&results)
	}
	if err != nil {
		return results, err
	}
	return results, err
}

func WriteRace(candidateRace ParsedRace) {
	file, _ := os.Create(FILEPATH)
	defer file.Close()
	json.NewEncoder(file).Encode(candidateRace)
	log.Printf("New race written - %v.", candidateRace)
}

func LoadRace() (storedRace ParsedRace) {
	data, err := os.ReadFile(FILEPATH)
	if err != nil {
		return
	}
	_ = json.Unmarshal(data, &storedRace)
	return
}

func SendTGMessage(text string) {
	tgToken, status := os.LookupEnv("TG_TOKEN")
	if !status {
		log.Printf("TG_TOKEN env is missing.")
		os.Exit(1)
	}

	notificationID, status := os.LookupEnv("NOTIFICATION_ID")
	if !status {
		log.Printf("NOTIFICATION_ID env is missing.")
		os.Exit(1)
	}
	notificationIDNum, _ := strconv.ParseInt(notificationID, 10, 64)
	bot, _ := t.NewBotAPI(tgToken)
	msg := t.NewMessage(notificationIDNum, text)

	var err error

	for {
		_, err = bot.Send(msg)
		if err == nil {
			break
		}
		log.Print(err)
		time.Sleep(30 * time.Second)
	}
}

func ParseRase(race Race) (parsedRace ParsedRace) {
	parsedRace.City = race.Circuit.Location.Locality
	parsedRace.Country = race.Circuit.Location.Country
	parsedRace.RaceDate, _ = time.Parse(DATE_FORMAT, CombineDate(race.Date, race.Time))
	parsedRace.QualDate, _ = time.Parse(DATE_FORMAT, CombineDate(race.Qualifying.Date, race.Qualifying.Time))
	if hasSprintSet(race) {
		parsedRace.ShootDate, _ = time.Parse(DATE_FORMAT, CombineDate(race.SprintQualifying.Date, race.SprintQualifying.Time))
		parsedRace.SprintDate, _ = time.Parse(DATE_FORMAT, CombineDate(race.Sprint.Date, race.Sprint.Time))
	}
	return
}

func CombineDate(date string, time string) string {
	return date + " " + time
}

func hasSprintSet(race Race) bool {
	zeroValue := reflect.Zero(reflect.TypeOf(race.Sprint)).Interface()
	return !reflect.DeepEqual(race.Sprint, zeroValue)
}

func isSameDateTime(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() &&
		t1.Month() == t2.Month() &&
		t1.Day() == t2.Day() &&
		t1.Hour() == t2.Hour() &&
		t1.Minute() == t2.Minute()
}
