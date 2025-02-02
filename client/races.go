package client

import (
    "fmt"
    "log"
    "time"
    "strconv"
    "reflect"
    "encoding/json"
    "github.com/enescakir/emoji"
)

type RacesResponse struct {
	MRData RacesMRData `json:"MRData"`
}

type RacesMRData struct {
	XMLNS     string    `json:"xmlns"`
	Series    string    `json:"series"`
	URL       string    `json:"url"`
	Limit     string    `json:"limit"`
	Offset    string    `json:"offset"`
	Total     string    `json:"total"`
	RaceTable RaceTable `json:"RaceTable"`
}

type RaceTable struct {
	Season string `json:"season"`
	Race  []Race  `json:"Races"`
}

type Race struct {
	Season           string  `json:"season"`
	Round            string  `json:"round"`
	URL              string  `json:"url"`
	RaceName         string  `json:"raceName"`
	Circuit          Circuit `json:"Circuit"`
	Date             string  `json:"date"`
	Time             string  `json:"time"`
	FirstPractice    Session `json:"FirstPractice"`
	SecondPractice   Session `json:"SecondPractice,omitempty"`
	ThirdPractice    Session `json:"ThirdPractice,omitempty"`
    Sprint           Session `json:"Sprint,omitempty"`
	Qualifying       Session `json:"Qualifying"`
    SprintQualifying Session `json:"SprintQualifying,omitempty"`
}

type Circuit struct {
	CircuitID   string   `json:"circuitId"`
	URL         string   `json:"url"`
	CircuitName string   `json:"circuitName"`
	Location    Location `json:"Location"`
}

type Location struct {
	Lat      string `json:"lat"`
	Long     string `json:"long"`
	Locality string `json:"locality"`
	Country  string `json:"country"`
}

type Session struct {
	Date string `json:"date"`
	Time string `json:"time"`
}


var countryFlagMap = map[string]string{
    "Japan": emoji.Parse(":jp:"),
    "Saudi Arabia": emoji.Parse(":saudi_arabia:"),
    "Azerbaijan": emoji.Parse(":azerbaijan:"),
    "USA": emoji.Parse(":us:"),
    "Monaco": emoji.Parse(":monaco:"),
    "Spain": emoji.Parse(":es:"),
    "Canada": emoji.Parse(":canada:"),
    "Austria": emoji.Parse(":austria:"),
    "UK": emoji.Parse(":gb:"),
    "Hungary": emoji.Parse(":hungary:"),
    "Belgium": emoji.Parse(":belgium:"),
    "Netherlands": emoji.Parse(":netherlands:"),
    "Italy": emoji.Parse(":it:"),
    "Singapore": emoji.Parse(":singapore:"),
    "Qatar": emoji.Parse(":qatar:"),
    "Mexico": emoji.Parse(":mexico:"),
    "Brazil": emoji.Parse(":brazil:"),
    "UAE": emoji.Parse(":united_arab_emirates:"),
    "United States": emoji.Parse(":us:"),
    "Bahrain": emoji.Parse(":bahrain:"),
    "Australia": emoji.Parse(":australia:"),
    "China": emoji.Parse(":flag-cn:"),
}


func RacesParseData(requestBody []byte) string {

    var response string
    var racesData RacesResponse
    var raceStartTime string
    var raceTime bool
    TBA := "\nTo be announced"
    err := json.Unmarshal(requestBody, &racesData)
    if err != nil {
        log.Println(err)
    }

	englishLocation, err := time.LoadLocation("Europe/London")
	if err != nil {
		fmt.Println("Error loading time zone:", err)
		return "Failed to fetch timezone"
	}

    for _, elem := range racesData.MRData.RaceTable.Race {
        raceTime = true
        if hasSprintSet(elem) {
            elem.Sprint.Time = ParseTime(elem.Date, elem.Sprint.Time, false)
            elem.SprintQualifying.Time = ParseTime(elem.Date, elem.SprintQualifying.Time, false)
        }

        if len(elem.Time) == 0 {
            elem.Time = "00:00:00Z"
            raceTime = false
        }
        raceStartTime = elem.Date + "T" + elem.Time
        raceStart, _ := time.Parse(time.RFC3339, raceStartTime)
        englishStartTime := raceStart.In(englishLocation)
        currentDate := time.Now()

        // Uncomment for testing
        //currentDate, _ := time.Parse("2006-01-02 15:04:05", "2024-12-23 05:00:00")

        elem.Time = ParseTime(elem.Date, elem.Time, false)
        elem.Qualifying.Time = ParseTime(elem.Date, elem.Qualifying.Time, false)

        flagEmoji, foundEmoji := countryFlagMap[elem.Circuit.Location.Country]

        if !foundEmoji {
            flagEmoji = "🏴‍☠️"
        }

        if currentDate.Before(englishStartTime) {
            // Hacky, maybe there is a way to make it better
            if raceTime == false {
                elem.Qualifying.Time = TBA
                elem.Time = TBA
            }
            response = fmt.Sprintf("Next race will take place in:\n%s %s %s\n",
                flagEmoji,
                elem.Circuit.Location.Country,
                elem.Circuit.Location.Locality,
            )
            if hasSprintSet(elem) {
                // Hacky, maybe there is a way to make it better
                if raceTime == false {
                    elem.SprintQualifying.Time = TBA
                    elem.Sprint.Time = TBA
                }
                response = response + fmt.Sprintf("\n🔫 Sprint Shootout - %s%s \n",
                    fmtDate(elem.SprintQualifying.Date),
                    elem.SprintQualifying.Time,
                )
                response = response + fmt.Sprintf("\n🏎 Sprint - %s%s \n",
                    fmtDate(elem.Sprint.Date),
                    elem.Sprint.Time,
                )
            }
            response = response + fmt.Sprintf("\n⏱ Qualifying - %s%s\n",
                fmtDate(elem.Qualifying.Date),
                elem.Qualifying.Time,
            )
            response = response + fmt.Sprintf("\n🏁 Race - %s%s",
                fmtDate(elem.Date),
                elem.Time,
            )
            break
        }
    }
    if len(response) == 0 {
        response = "No more races this year 😭"
    }
    return response
}


func ParseTime (rawDate string, rawTime string, fixTime bool) string {
	spanishLocation, err := time.LoadLocation("Europe/Madrid")
	if err != nil {
		fmt.Println("Error loading time zone:", err)
		return "Failed to fetch timezone"
	}

	englishLocation, err := time.LoadLocation("Europe/London")
	if err != nil {
		fmt.Println("Error loading time zone:", err)
		return "Failed to fetch timezone"
	}

    var processedTime string
    // Adjust date and time to follow this format for parsing - 2025-06-01T15:04:05Z
    dateTime := rawDate + "T" + rawTime
    adjustedTime, _ := time.Parse(time.RFC3339, dateTime)
    // Probably not needed anymore, TODO - verify later this year
    // Adjust time for Sprint Shootouts because API is reporting it with incorrect time
    /*if fixTime == true {
        adjustedTime = adjustedTime.Add(-30 * time.Minute)
    }*/
    
    // Both these valuse come with following format - 2024-10-20 23:00:00 +0100 BST
    spanishTime := adjustedTime.In(spanishLocation)
    englishTime := adjustedTime.In(englishLocation)

    processedTime = fmt.Sprintf("\n%s%s\n%s%s",
        countryFlagMap["UK"],
        englishTime.Format("15:04"),
        countryFlagMap["Spain"],
        spanishTime.Format("15:04"),
    )

    return processedTime
}


func rmSeconds(longTime string) string {
    return longTime[:len(longTime)-3]
}


func fmtDate(longDate string) string {
    longDate = longDate[5:]
    monthValue, _ := strconv.Atoi(longDate[:2])
    return time.Month(monthValue).String() + " " + longDate[3:]
}


func hasSprintSet(race Race) bool {
    zeroValue := reflect.Zero(reflect.TypeOf(race.Sprint)).Interface()
    return !reflect.DeepEqual(race.Sprint, zeroValue)
}
