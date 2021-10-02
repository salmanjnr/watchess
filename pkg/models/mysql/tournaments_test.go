package mysql

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"watchess.org/watchess/pkg/models"
)

type tournamentSlice []*models.Tournament

// For debugging purposes
func (ts tournamentSlice) String() string {
	s := "["
	for i, tournement := range ts {
		if i > 0 {
			s += ", "
		}
		s += fmt.Sprintf("%v", tournement)
	}
	return s + "]"
}

func TestTournamentModelGet(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	tests := []struct {
		name           string
		tournamentID   int
		wantTournament *models.Tournament
		wantError      error
	}{
		{
			name:         "Valid ID",
			tournamentID: 3,
			wantTournament: &models.Tournament{
				ID:               3,
				Name:             "Norway Chess",
				ShortDescription: "Watch Magnus",
				LongDescription:  "Watch Magnus and Nepo",
				HasStandings:     true,
				StartDate:        time.Date(2021, 9, 7, 17, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2021, 9, 17, 22, 0, 0, 0, time.UTC),
				IsLive:           false,
				OwnerID:          1,
			},
			wantError: nil,
		},
		{
			name:           "Zero ID",
			tournamentID:   0,
			wantTournament: nil,
			wantError:      models.ErrNoRecord,
		},
		{
			name:           "Non-existent ID",
			tournamentID:   6927,
			wantTournament: nil,
			wantError:      models.ErrNoRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			m := TournamentModel{db}

			tournament, err := m.Get(tt.tournamentID)

			if err != tt.wantError {
				t.Errorf("want %v; got %s", tt.wantError, err)
			}

			if !reflect.DeepEqual(tournament, tt.wantTournament) {
				t.Errorf("want %v; got %v", tt.wantTournament, tournament)
			}
		})
	}
}

func TestTournamentModelScan(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	tests := []struct {
		name            string
		function        func(TournamentModel) ([]*models.Tournament, error)
		wantTournaments []*models.Tournament
		wantError       error
	}{
		{
			name: "LatestActive",
			function: func(m TournamentModel) ([]*models.Tournament, error) {
				return m.LatestActive(5)
			},
			wantTournaments: []*models.Tournament{
				{
					ID:               1,
					Name:             "Happening Now",
					ShortDescription: "Watch now",
					LongDescription:  "Watch the best in the world lose",
					HasStandings:     false,
					StartDate:        time.Date(2021, 3, 7, 12, 0, 0, 0, time.UTC),
					EndDate:          time.Date(2321, 3, 17, 20, 0, 0, 0, time.UTC),
					IsLive:           true,
					OwnerID:          1,
				},
				{
					ID:               2,
					Name:             "Also Now",
					ShortDescription: "Watch now",
					LongDescription:  "Watch the best in the world lose",
					HasStandings:     false,
					StartDate:        time.Date(2021, 5, 7, 12, 0, 0, 0, time.UTC),
					EndDate:          time.Date(2321, 3, 17, 20, 0, 0, 0, time.UTC),
					IsLive:           true,
					OwnerID:          1,
				},
			},
			wantError: nil,
		},
		{
			name: "LatestFinished",
			function: func(m TournamentModel) ([]*models.Tournament, error) {
				return m.LatestFinished(5)
			},
			wantTournaments: []*models.Tournament{
				{
					ID:               3,
					Name:             "Norway Chess",
					ShortDescription: "Watch Magnus",
					LongDescription:  "Watch Magnus and Nepo",
					HasStandings:     true,
					StartDate:        time.Date(2021, 9, 7, 17, 0, 0, 0, time.UTC),
					EndDate:          time.Date(2021, 9, 17, 22, 0, 0, 0, time.UTC),
					IsLive:           false,
					OwnerID:          1,
				},
				{
					ID:               4,
					Name:             "Candidates",
					ShortDescription: "Watch good players",
					LongDescription:  "Watch good players fight for the chance to play against Magnus",
					HasStandings:     true,
					StartDate:        time.Date(2021, 3, 7, 11, 0, 0, 0, time.UTC),
					EndDate:          time.Date(2021, 3, 17, 12, 0, 0, 0, time.UTC),
					IsLive:           false,
					OwnerID:          1,
				},
			},
			wantError: nil,
		},
		{
			name: "Upcoming",
			function: func(m TournamentModel) ([]*models.Tournament, error) {
				return m.Upcoming(5)
			},
			wantTournaments: []*models.Tournament{
				{
					ID:               5,
					Name:             "Candidates Far Far Far",
					ShortDescription: "Watch good players probably",
					LongDescription:  "Watch good players probably fight for the chance to play against Magnus",
					HasStandings:     true,
					StartDate:        time.Date(2321, 3, 7, 7, 0, 0, 0, time.UTC),
					EndDate:          time.Date(2321, 3, 17, 10, 0, 0, 0, time.UTC),
					IsLive:           false,
					OwnerID:          1,
				},
				{
					ID:               6,
					Name:             "World Championship Far Far",
					ShortDescription: "Watch Magnus vs someone",
					LongDescription:  "Watch Magnus vs someone that might be Alireze",
					HasStandings:     false,
					StartDate:        time.Date(2323, 3, 7, 7, 0, 0, 0, time.UTC),
					EndDate:          time.Date(2324, 3, 17, 10, 0, 0, 0, time.UTC),
					IsLive:           false,
					OwnerID:          1,
				},
			},
			wantError: nil,
		},
		{
			name: "Upcoming Limited",
			function: func(m TournamentModel) ([]*models.Tournament, error) {
				return m.Upcoming(1)
			},
			wantTournaments: []*models.Tournament{
				{
					ID:               5,
					Name:             "Candidates Far Far Far",
					ShortDescription: "Watch good players probably",
					LongDescription:  "Watch good players probably fight for the chance to play against Magnus",
					HasStandings:     true,
					StartDate:        time.Date(2321, 3, 7, 7, 0, 0, 0, time.UTC),
					EndDate:          time.Date(2321, 3, 17, 10, 0, 0, 0, time.UTC),
					IsLive:           false,
					OwnerID:          1,
				},
			},
			wantError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			m := TournamentModel{db}

			tournament, err := tt.function(m)

			if err != tt.wantError {
				t.Errorf("want %v; got %s", tt.wantError, err)
			}

			if !reflect.DeepEqual(tournament, tt.wantTournaments) {
				t.Errorf("want %v; got %v", tournamentSlice(tt.wantTournaments), tournamentSlice(tournament))
			}
		})
	}
}

func TestTournamentModelInsert(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	tests := []struct {
		name       string
		tournament *models.Tournament
		wantError  error
	}{
		{
			name: "Valid",
			tournament: &models.Tournament{
				ID:               7,
				Name:             "Test",
				ShortDescription: "Test Description",
				LongDescription:  "Test Long Description",
				HasStandings:     true,
				StartDate:        time.Date(2021, 9, 7, 17, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2021, 9, 17, 22, 0, 0, 0, time.UTC),
				IsLive:           false,
				OwnerID:          1,
			},
			wantError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			m := TournamentModel{db}

			tournamentID, err := m.Insert(tt.tournament.Name, tt.tournament.ShortDescription, tt.tournament.LongDescription, tt.tournament.HasStandings, tt.tournament.StartDate, tt.tournament.EndDate, tt.tournament.IsLive, tt.tournament.OwnerID)

			if err != tt.wantError {
				t.Errorf("want %v; got %s", tt.wantError, err)
			}

			tournament, err := m.Get(tournamentID)

			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(tournament, tt.tournament) {
				t.Errorf("want %v; got %v", tt.tournament, tournament)
			}
		})
	}
}
