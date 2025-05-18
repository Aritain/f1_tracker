package client

import (
	"fmt"
	"time"

	"github.com/enescakir/emoji"
)

var countryFlagMap = map[string]string{
	"Japan":         emoji.Parse(":jp:"),
	"Saudi Arabia":  emoji.Parse(":saudi_arabia:"),
	"Azerbaijan":    emoji.Parse(":azerbaijan:"),
	"USA":           emoji.Parse(":us:"),
	"Monaco":        emoji.Parse(":monaco:"),
	"Spain":         emoji.Parse(":es:"),
	"Canada":        emoji.Parse(":canada:"),
	"Austria":       emoji.Parse(":austria:"),
	"UK":            emoji.Parse(":gb:"),
	"Hungary":       emoji.Parse(":hungary:"),
	"Belgium":       emoji.Parse(":belgium:"),
	"Netherlands":   emoji.Parse(":netherlands:"),
	"Italy":         emoji.Parse(":it:"),
	"Singapore":     emoji.Parse(":singapore:"),
	"Qatar":         emoji.Parse(":qatar:"),
	"Mexico":        emoji.Parse(":mexico:"),
	"Brazil":        emoji.Parse(":brazil:"),
	"UAE":           emoji.Parse(":united_arab_emirates:"),
	"United States": emoji.Parse(":us:"),
	"Bahrain":       emoji.Parse(":bahrain:"),
	"Australia":     emoji.Parse(":australia:"),
	"China":         emoji.Parse(":flag-cn:"),
}

// mode = query, sprint, race
func PrepRaceMessage(race ParsedRace, mode string) (text string) {
	var dummyTime time.Time
	var ukTime time.Time
	var spTime time.Time

	// Prepare first line
	var topLine string
	switch mode {
	case "sprint":
		topLine = "🏎 Sprint today in:\n"
	case "race":
		topLine = "🏁 Race today in:\n"
	}
	text += topLine

	// Prepare flag + location
	flagEmoji, foundEmoji := countryFlagMap[race.Country]

	if !foundEmoji {
		flagEmoji = "🏴‍☠️"
	}
	locationLine := fmt.Sprintf("%s %s %s\n", flagEmoji, race.Country, race.City)
	text += locationLine

	// Prep sprint time
	if mode == "sprint" {
		ukTime, spTime = LocalTime(race.SprintDate)
	}

	// Prep race time
	if mode == "race" {
		ukTime, spTime = LocalTime(race.RaceDate)
	}

	// Prepare sprint/race msg
	if mode == "sprint" || mode == "race" {
		ukHours := ukTime.Format("15:04")
		spHours := spTime.Format("15:04")
		text += fmt.Sprintf("🇬🇧 %s\n🇪🇸 %s", ukHours, spHours)
		return
	}

	// Query message block

	if race.RaceDate == dummyTime {
		text = "No more races this year 😭"
		return
	}

	if race.ShootDate != dummyTime {
		text += "\n🔫 Shootout - "
		ukTime, spTime = LocalTime(race.ShootDate)
		text += ParseDateTime(ukTime, spTime)

		text += "\n🏎 Sprint - "
		ukTime, spTime = LocalTime(race.SprintDate)
		text += ParseDateTime(ukTime, spTime)
	}

	text += "\n⏱ Qualifying - "
	ukTime, spTime = LocalTime(race.QualDate)
	text += ParseDateTime(ukTime, spTime)

	text += "\n🏁 Race - "
	ukTime, spTime = LocalTime(race.RaceDate)
	text += ParseDateTime(ukTime, spTime)

	return
}

func ParseDateTime(ukTime time.Time, spTime time.Time) (text string) {
	ukHours := ukTime.Format("15:04")
	spHours := spTime.Format("15:04")
	date := ukTime.Format("January 2")
	text = fmt.Sprintf("%s\n🇬🇧 %s\n🇪🇸 %s\n", date, ukHours, spHours)
	return
}

func LocalTime(date time.Time) (ukTime time.Time, spTime time.Time) {
	englishLocation, _ := time.LoadLocation("Europe/London")
	spanishLocation, _ := time.LoadLocation("Europe/Madrid")
	ukTime = date.In(englishLocation)
	spTime = date.In(spanishLocation)
	return
}
