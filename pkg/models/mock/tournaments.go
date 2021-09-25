package mock

import (
	"time"

	"watchess.org/watchess/pkg/models"
)

var mockTournament1 = &models.Tournament{
	ID:               1,
	Name:             "Norway Chess",
	ShortDescription: "Magnus returns to his home country to compete against many 2700+",
	LongDescription:  "Magnus returns to his home country to compete against many 2700+, including Nepo, his challenger for the WC title.",
	HasStandings:     true,
	StartDate:        time.Date(2021, time.Month(9), 7, 17, 0, 0, 0, time.UTC),
	EndDate:          time.Date(2021, time.Month(9), 17, 22, 0, 0, 0, time.UTC),
	IsLive:           false,
}

var mockTournament2 = &models.Tournament{
	ID:               2,
	Name:             "Future World Chess Championship",
	ShortDescription: "Magnus vs challenger",
	LongDescription:  "Magnus vs challenger (possibly Alireze)",
	HasStandings:     false,
	StartDate:        time.Date(2031, time.Month(9), 7, 17, 0, 0, 0, time.UTC),
	EndDate:          time.Date(2031, time.Month(9), 17, 22, 0, 0, 0, time.UTC),
	IsLive:           false,
}

var mockTournament3 = &models.Tournament{
	ID:               3,
	Name:             "Happening Now",
	ShortDescription: "This tournament is live",
	LongDescription:  "",
	HasStandings:     true,
	StartDate:        time.Date(2021, time.Month(9), 7, 17, 0, 0, 0, time.UTC),
	EndDate:          time.Date(2031, time.Month(9), 17, 22, 0, 0, 0, time.UTC),
	IsLive:           true,
}

type TournamentModel struct{}

func (m *TournamentModel) Insert(name, shortDescription, longDescription string, hasStandings bool, startDate, endDate time.Time, isLive bool) (int, error) {
	return 1, nil
}

func (m *TournamentModel) Get(id int) (*models.Tournament, error) {
	switch id {
	case 1:
		return mockTournament1, nil
	case 2:
		return mockTournament2, nil
	case 3:
		return mockTournament3, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *TournamentModel) LatestActive(limit int) ([]*models.Tournament, error) {
	return []*models.Tournament{mockTournament3}, nil
}

func (m *TournamentModel) LatestFinished(limit int) ([]*models.Tournament, error) {
	return []*models.Tournament{mockTournament1}, nil
}

func (m *TournamentModel) Upcoming(limit int) ([]*models.Tournament, error) {
	return []*models.Tournament{mockTournament2}, nil
}
