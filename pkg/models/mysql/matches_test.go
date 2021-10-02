package mysql

import (
	"reflect"
	"testing"

	"watchess.org/watchess/pkg/models"
)

func TestMatchesModelInsert(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	tests := []struct {
		name       string
		match *models.Match
		wantError  error
	}{
		{
			name: "Valid",
			match: &models.Match{
				Side1: "Foo1",
				Side2: "Foo2",
				RoundID: 2,
			},
			wantError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			m := MatchModel{db}

			matchID, err := m.Insert(tt.match.Side1, tt.match.Side2, tt.match.RoundID)

			if err != tt.wantError {
				t.Errorf("want %v; got %s", tt.wantError, err)
			}

			match, err := m.Get(matchID)

			if err != nil {
				t.Fatal(err)
			}

			tt.match.ID = matchID

			if !reflect.DeepEqual(match, tt.match) {
				t.Errorf("want %v; got %v", tt.match, match)
			}
		})
	}
}

func TestMatchesModelGet(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	tests := []struct {
		name           string
		matchID   int
		wantMatch *models.Match
		wantError      error
	}{
		{
			name:         "Valid ID",
			matchID: 3,
			wantMatch: &models.Match{
				ID: 3,
				Side1: "Alireza",
				Side2: "Richard",
				RoundID: 2,
			},
			wantError: nil,
		},
		{
			name:           "Zero ID",
			matchID:   0,
			wantMatch: nil,
			wantError:      models.ErrNoRecord,
		},
		{
			name:           "Non-existent ID",
			matchID:   6927,
			wantMatch: nil,
			wantError:      models.ErrNoRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			m := MatchModel{db}

			match, err := m.Get(tt.matchID)

			if err != tt.wantError {
				t.Errorf("want %v; got %s", tt.wantError, err)
			}

			if !reflect.DeepEqual(match, tt.wantMatch) {
				t.Errorf("want %v; got %v", tt.wantMatch, match)
			}
		})
	}
}
