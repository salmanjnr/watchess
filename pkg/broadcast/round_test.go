package broadcast

import (
	"strings"
	"testing"
	"time"

	"go.uber.org/goleak"
)

func TestRoundRun(t *testing.T) {
	defer goleak.VerifyNone(t)
	tests := []struct {
		name string
		pgnDir string
		wantMsgs []string
	}{
		{
			"valid",
			"valid",
			[]string{
				"2 matches and 5 games detected",
				"2 games updated",
				"2 games updated",
				"3 games updated",
				"1 game updated",
				"Closing",
			},
		},
		{
			"missing tags",
			"invalid_pgn/missing_tags",
			[]string{
				"Game 1 missing required tag: White",
				"Closing",
			},
		},
		{
			"invalid format",
			"invalid_pgn/invalid_format",
			[]string{
				"pgn decode error",
				"Closing",
			},
		},
		{
			"invalid updates",
			"invalid_updates",
			[]string{
				"2 matches and 5 games detected",
				"2 games updated",
				"2 games updated",
				"3 games updated",
				"Cannot receive update for finished game number 1",
				"Closing",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts, err := newPGNServer(t, tt.pgnDir)
			if err != nil {
				panic(err)
			}
			defer ts.Close()
			round := newTestRound(ts)
			timeout := time.NewTicker(time.Duration(2) * time.Second)
			defer timeout.Stop()

			go round.Run()
			logChan := round.LogChan
			i := 0
			for {
				if i == len(tt.wantMsgs) {
					break
				}
				select {
				case update, ok := <- logChan:
					if !ok {
						logChan = nil
						break
					}
					var msg string
					switch uType := update.(type) {
					case string:
						msg = update.(string)
					case error:
						msg = update.(error).Error()
					default:
						t.Errorf("Received log message has unsupported type %v", uType)
						i += 1
						break
					}

					if !strings.Contains(msg, tt.wantMsgs[i]) {
						t.Errorf("Error in message number %d. Want: %v. Received: %v", i, tt.wantMsgs[i], msg)
					}
					i += 1
				case <-timeout.C:
					t.Error("Timedout before receiving wanted messages")
				}
				}
		})	
	}
}
