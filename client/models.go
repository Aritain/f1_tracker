package client

import "time"

// Race
type RacesResponse struct {
	MRData RacesMRData `json:"MRData"`
}

type RacesMRData struct {
	RaceTable RaceTable `json:"RaceTable"`
}

type RaceTable struct {
	Race []Race `json:"Races"`
}

type Race struct {
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

// Team
type TeamsResponse struct {
	MRData TeamsMRData `json:"MRData"`
}

type TeamsMRData struct {
	StandingsTable TeamStandingsTable `json:"StandingsTable"`
}

type TeamStandingsTable struct {
	StandingsLists []TeamStandingsLists `json:"StandingsLists"`
}

type TeamStandingsLists struct {
	ConstructorStandings []ConstructorStandings `json:"ConstructorStandings"`
}

type ConstructorStandings struct {
	Points      string          `json:"points"`
	Constructor TeamConstructor `json:"Constructor"`
}

type TeamConstructor struct {
	Name string `json:"name"`
}

// Driver
type DriversResponse struct {
	MRData DriverMRData `json:"MRData"`
}

type DriverMRData struct {
	StandingsTable DriverStandingsTable `json:"StandingsTable"`
}

type DriverStandingsTable struct {
	StandingsList []DriverStandingsList `json:"StandingsLists"`
}

type DriverStandingsList struct {
	DriverStanding []DriverStanding `json:"DriverStandings"`
}

type DriverStanding struct {
	Points string `json:"points"`
	Driver Driver `json:"Driver"`
}

type Driver struct {
	Code string `json:"code"`
}

// Stored race
type StoredRace struct {
	Date time.Time `json:"date"`
}
