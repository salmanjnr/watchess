package mysql

import (
	"testing"

	"watchess.org/watchess/pkg/models"
)

func TestRoundModelDelete(t *testing.T) {
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

			r := RoundModel{db}

			err := r.Delete(tt.roundID)
			if err != tt.wantError {
				t.Errorf("want %v; got %s", tt.wantError, err)
			}

			if err != nil {
				return
			}

			_, err2 := r.Get(tt.roundID)
			if err2 != models.ErrNoRecord {
				t.Errorf("Unexpected error when fetching game %v", err2)
			}
		})
	}
}
