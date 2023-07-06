package client

import (
    "log"
    "os"
    "time"
    "strconv"
    "reflect"
    "encoding/xml"
    "github.com/enescakir/emoji"
)

type RacesMRData struct {
    XMLName   xml.Name  `xml:"MRData"`
    RaceTable RaceTable `xml:"RaceTable"`
}

type RaceTable struct {
    XMLName  xml.Name `xml:"RaceTable"`
    Race     []Race   `xml:"Race"`
}

type Race struct {
    XMLName        xml.Name `xml:"Race"`
    RaceName       string   `xml:"RaceName"`
    Circuit        Circuit  `xml:"Circuit"`
    Date           string   `xml:"Date"`
    Time           string   `xml:"Time"`
    FirstPractice  PreRace  `xml:"FirstPractice"`
    SecondPractice PreRace  `xml:"SecondPractice"`
    ThirdPractice  PreRace  `xml:"ThirdPractice"`
    Sprint         PreRace  `xml:"Sprint,omitempty"`
    Qualifying     PreRace  `xml:"Qualifying"`
}

type Circuit struct {
    CircuitId   string   `xml:"circuitId"`
    URL         string   `xml:"url"`
    CircuitName string   `xml:"CircuitName"`
    Location    Location `xml:"Location"`
}

type Location struct {
    Locality string `xml:"Locality"`
    Country  string `xml:"Country"`
}

type PreRace struct {
    Date string `xml:"Date"`
    Time string `xml:"Time"`
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
}


func RacesParseData(requestBody []byte) string {

    _, status := os.LookupEnv("TZ_OFFSET")
    if status == false {
        log.Printf("TZ_OFFSET env is missing.")
        os.Exit(1)
    }
    tzOffset, _ := strconv.Atoi(os.Getenv("TZ_OFFSET"))

    var response string
    var racesData RacesMRData
    var raceEndTime string
    err := xml.Unmarshal(requestBody, &racesData)
    if err != nil {
        log.Println(err)
    }
    
    for _, elem := range racesData.RaceTable.Race {
        if hasSprintSet(elem) {
            elem.Sprint.Time = ParseTime(elem.Sprint.Time, false, tzOffset)
            elem.SecondPractice.Time = ParseTime(elem.SecondPractice.Time, true, tzOffset)
        }

        // Races are usually two hours longs
        raceEndTime = ParseTime(elem.Time, false, tzOffset+2)
        raceEndTime = elem.Date + " " +  raceEndTime
        raceEnd, _ := time.Parse("2006-01-02 15:04:05", raceEndTime)
        currentDate := time.Now()

        elem.Time = ParseTime(elem.Time, false, tzOffset)
        elem.Qualifying.Time = ParseTime(elem.Qualifying.Time, false, tzOffset)

        flagEmoji, foundEmoji := countryFlagMap[elem.Circuit.Location.Country]

        if !foundEmoji {
            flagEmoji = "🏴‍☠️"
        }

        if currentDate.Before(raceEnd) {
            response = "Next race will take place in:\n" + 
                flagEmoji + " " + elem.Circuit.Location.Country + 
                " " + elem.Circuit.Location.Locality + "\n⏱ Qualifying: " + 
                elem.Qualifying.Date + " " + elem.Qualifying.Time + "\n"
            if hasSprintSet(elem) {
                response = response + "🔫 Sprint Shootout: " + elem.SecondPractice.Date + 
                    " " + elem.SecondPractice.Time + "\n"
                response = response + "🏎 Sprint: " + elem.Sprint.Date + " " + elem.Sprint.Time + "\n"
            }
            response = response + "🏁 Race: " + elem.Date + " " + elem.Time
            break
        }
    }
    return response
}

func ParseTime (rawTime string, fixTime bool, tzOffset int) string {

    adjustedTime, _ := time.Parse("15:04:05Z", rawTime)
    adjustedTime = adjustedTime.Add(time.Duration(tzOffset) * time.Hour)
    // Adjust time for Sprint Shootouts because API is reporting it with incorrect time
    if fixTime == true {
        adjustedTime = adjustedTime.Add(-30 * time.Minute)
    }
    return adjustedTime.Format("15:04:05")

}

func hasSprintSet(race Race) bool {
    zeroValue := reflect.Zero(reflect.TypeOf(race.Sprint)).Interface()
    return !reflect.DeepEqual(race.Sprint, zeroValue)
}
