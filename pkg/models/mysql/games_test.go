package mysql

import (
	"fmt"
	"reflect"
	"testing"

	"watchess.org/watchess/pkg/models"
)

type gameSlice []*models.Game

// For debugging purposes
func (gs gameSlice) String() string {
	s := "["
	for i, game := range gs {
		if i > 0 {
			s += ", "
		}
		s += fmt.Sprintf("%v", game)
	}
	return s + "]"
}

func getGameResult(res string) models.GameResult {
	result, _ := models.GetGameResult(res)
	return *result
}

func TestGameModelInsert(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	tests := []struct {
		name      string
		game      *models.Game
		wantError error
	}{
		{
			name: "ValidWin",
			game: &models.Game{
				White:          "Foo1",
				Black:          "Foo2",
				Result:         getGameResult("1-0"),
				WhiteMatchSide: "Foo1",
				BlackMatchSide: "Foo2",
				PGN:            "Test",
				MatchID:        1,
				RoundID:        1,
			},
			wantError: nil,
		},
		{
			name: "ValidDraw",
			game: &models.Game{
				White:          "Foo1",
				Black:          "Foo2",
				Result:         getGameResult("1/2-1/2"),
				WhiteMatchSide: "Foo1",
				BlackMatchSide: "Foo2",
				PGN:            "Test",
				MatchID:        1,
				RoundID:        1,
			},
			wantError: nil,
		},
		{
			name: "ValidLose",
			game: &models.Game{
				White:          "Foo1",
				Black:          "Foo2",
				Result:         getGameResult("0-1"),
				WhiteMatchSide: "Foo1",
				BlackMatchSide: "Foo2",
				PGN:            "Test",
				MatchID:        1,
				RoundID:        1,
			},
			wantError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			m := GameModel{db}

			gameID, err := m.Insert(tt.game.White, tt.game.Black, tt.game.Result, tt.game.WhiteMatchSide, tt.game.BlackMatchSide, tt.game.PGN, tt.game.MatchID, tt.game.RoundID)

			if err != tt.wantError {
				t.Errorf("want %v; got %v", tt.wantError, err)
			}

			game, err := m.Get(gameID)

			if err != nil {
				t.Fatal(err)
			}

			tt.game.ID = gameID

			if !reflect.DeepEqual(game, tt.game) {
				t.Errorf("want %v; got %v", tt.game, game)
			}
		})
	}
}

func TestGameModelGet(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	tests := []struct {
		name      string
		gameID    int
		wantGame  *models.Game
		wantError error
	}{
		{
			name:   "Valid ID",
			gameID: 4,
			wantGame: &models.Game{
				ID:             4,
				White:          "Alireza",
				Black:          "Richard",
				Result:         getGameResult("0-1"),
				WhiteMatchSide: "Alireza",
				BlackMatchSide: "Richard",
				PGN:            "Test",
				MatchID:        3,
				RoundID:        2,
			},
			wantError: nil,
		},
		{
			name:      "Zero ID",
			gameID:    0,
			wantGame:  nil,
			wantError: models.ErrNoRecord,
		},
		{
			name:      "Non-existent ID",
			gameID:    6927,
			wantGame:  nil,
			wantError: models.ErrNoRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			g := GameModel{db}

			game, err := g.Get(tt.gameID)

			if err != tt.wantError {
				t.Errorf("want %v; got %s", tt.wantError, err)
			}

			if !reflect.DeepEqual(game, tt.wantGame) {
				t.Errorf("want %v; got %v", tt.wantGame, game)
			}
		})
	}
}

func TestGameModelGetByRound(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	tests := []struct {
		name      string
		roundID   int
		wantGames []*models.Game
		wantError error
	}{
		{
			name:    "Valid",
			roundID: 1,
			wantGames: []*models.Game{
				{
					ID:             1,
					White:          "Magnus Carlsen",
					Black:          "Ian Nepo",
					Result:         getGameResult("1-0"),
					WhiteMatchSide: "Magnus Carlsen",
					BlackMatchSide: "Ian Nepo",
					PGN:            "Test",
					MatchID:        1,
					RoundID:        1,
				},
				{
					ID:             2,
					White:          "Ian Nepo",
					Black:          "Magnus Carlsen",
					Result:         getGameResult("1/2-1/2"),
					WhiteMatchSide: "Ian Nepo",
					BlackMatchSide: "Magnus Carlsen",
					PGN:            "Test",
					MatchID:        1,
					RoundID:        1,
				},
				{
					ID:             3,
					White:          "Magnus Carlsen",
					Black:          "Ian Nepo",
					Result:         getGameResult("0-1"),
					WhiteMatchSide: "Magnus Carlsen",
					BlackMatchSide: "Ian Nepo",
					PGN:            "Test",
					MatchID:        1,
					RoundID:        1,
				},
			},
			wantError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			g := GameModel{db}

			games, err := g.GetByRound(tt.roundID)

			if err != tt.wantError {
				t.Errorf("want %v; got %s", tt.wantError, err)
			}

			if !reflect.DeepEqual(games, tt.wantGames) {
				t.Errorf("want %v; got %v", gameSlice(tt.wantGames), gameSlice(games))
			}
		})
	}

}

func TestGameModelUpdate(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	getStringRef := func(s string) *string {
		return &s
	}

	tests := []struct {
		name         string
		gameID       int
		updateValues struct {
			pgn    *string
			result *string
		}
		wantGame  *models.Game
		wantError error
	}{
		{
			name:   "Valid PGN",
			gameID: 4,
			updateValues: struct {
				pgn    *string
				result *string
			}{
				pgn: getStringRef("Hello"),
			},
			wantGame: &models.Game{
				ID:             4,
				White:          "Alireza",
				Black:          "Richard",
				Result:         getGameResult("0-1"),
				WhiteMatchSide: "Alireza",
				BlackMatchSide: "Richard",
				PGN:            "Hello",
				MatchID:        3,
				RoundID:        2,
			},
			wantError: nil,
		},
		{
			name:   "Valid Result",
			gameID: 4,
			updateValues: struct {
				pgn    *string
				result *string
			}{
				result: getStringRef("1-0"),
			},
			wantGame: &models.Game{
				ID:             4,
				White:          "Alireza",
				Black:          "Richard",
				Result:         getGameResult("1-0"),
				WhiteMatchSide: "Alireza",
				BlackMatchSide: "Richard",
				PGN:            "Test",
				MatchID:        3,
				RoundID:        2,
			},
			wantError: nil,
		},
		{
			name:   "Valid Both",
			gameID: 4,
			updateValues: struct {
				pgn    *string
				result *string
			}{
				pgn:    getStringRef("Hello"),
				result: getStringRef("1-0"),
			},
			wantGame: &models.Game{
				ID:             4,
				White:          "Alireza",
				Black:          "Richard",
				Result:         getGameResult("1-0"),
				WhiteMatchSide: "Alireza",
				BlackMatchSide: "Richard",
				PGN:            "Hello",
				MatchID:        3,
				RoundID:        2,
			},
			wantError: nil,
		},
		{
			name:   "Zero ID",
			gameID: 0,
			updateValues: struct {
				pgn    *string
				result *string
			}{
				pgn:    getStringRef("Hello"),
				result: getStringRef("1-0"),
			},
			wantGame:  nil,
			wantError: models.ErrNoRecord,
		},
		{
			name:   "Non-existent ID",
			gameID: 6927,
			updateValues: struct {
				pgn    *string
				result *string
			}{
				pgn:    getStringRef("Hello"),
				result: getStringRef("1-0"),
			},
			wantGame:  nil,
			wantError: models.ErrNoRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			g := GameModel{db}

			err := g.Update(tt.gameID, tt.updateValues.pgn, tt.updateValues.result)
			if err != tt.wantError {
				t.Errorf("want %v; got %s", tt.wantError, err)
			}

			if err != nil {
				return
			}

			game, err2 := g.Get(tt.gameID)

			if err2 != nil {
				t.Errorf("Unexpected error when fetching game %v", err2)
			}

			if !reflect.DeepEqual(game, tt.wantGame) {
				t.Errorf("want %v; got %v", tt.wantGame, game)
			}
		})
	}
}

func TestGameModelDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	tests := []struct {
		name      string
		gameID    int
		wantError error
	}{
		{
			name:      "Valid",
			gameID:    4,
			wantError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			g := GameModel{db}

			err := g.Delete(tt.gameID)
			if err != tt.wantError {
				t.Errorf("want %v; got %s", tt.wantError, err)
			}

			if err != nil {
				return
			}

			_, err2 := g.Get(tt.gameID)
			if err2 != models.ErrNoRecord {
				t.Errorf("Unexpected error when fetching game %v", err2)
			}
		})
	}
}

func TestGameModelDeleteByRound(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	tests := []struct {
		name      string
		roundID   int
		wantError error
	}{
		{
			name:      "Valid",
			roundID:   1,
			wantError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			g := GameModel{db}

			err := g.DeleteByRound(tt.roundID)
			if err != tt.wantError {
				t.Errorf("want %v; got %s", tt.wantError, err)
			}

			if err != nil {
				return
			}

			games, err2 := g.GetByRound(tt.roundID)
			if err2 != nil {
				t.Errorf("Unexpected error when fetching game %v", err2)
			}
			if len(games) != 0 {
				t.Errorf("Unexpected number of games. want 0; got %d", len(games))
			}
		})
	}
}
