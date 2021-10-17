package broadcast

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"go.uber.org/goleak"
)

func TestRoundRun(t *testing.T) {
	defer goleak.VerifyNone(t)
	tests := []struct {
		name     string
		pgnDir   string
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
				case update, ok := <-logChan:
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
					return
				}
			}
		})
	}
}

func TestRoundClient(t *testing.T) {
	defer goleak.VerifyNone(t)
	tests := []struct {
		name     string
		pgnDir   string
		wantMsgs []GameUpdate
	}{
		{
			"valid",
			"valid",
			[]GameUpdate{
				{
					GameID: 1,
					FEN: func() *string {
						s := "8/8/6p1/1p3p2/5P1q/1k2KP2/8/8 b - - 1 51"
						return &s
					}(),
				},
				{
					GameID: 4,
					FEN: func() *string {
						s := "2nq1r1k/1pb1R1pp/p1p1Qn2/P7/1P1P4/2P3P1/2B2P1P/4R1K1 w - - 3 36"
						return &s
					}(),
				},
				{
					GameID: 1,
					FEN: func() *string {
						s := "8/8/6p1/3q1p2/2k2P2/5P2/5K2/1q6 w - - 4 59"
						return &s
					}(),
				},
				{
					GameID: 4,
					FEN: func() *string {
						s := "3qr2k/1pb3pp/p1p2n2/P7/1P1P4/1QP3P1/2B2P1P/4R1K1 w - - 3 41"
						return &s
					}(),
				},
				{
					GameID: 1,
					FEN: func() *string {
						s := "8/8/6p1/5p2/2k2q2/5P1K/2q5/8 w - - 2 62"
						return &s
					}(),
					Result: func() *string {
						s := "1/2-1/2"
						return &s
					}(),
				},
				{
					GameID: 2,
					FEN: func() *string {
						s := "r6k/p5pp/2pb4/1p1p2q1/8/1BQ5/PPP3PP/5RK1 b - - 2 22"
						return &s
					}(),
				},
				{
					GameID: 4,
					FEN: func() *string {
						s := "3q3k/1pb3pp/p1p2n2/P7/1P1P4/1QP3P1/2B1rP1P/6K1 w - - 0 42"
						return &s
					}(),
					Result: func() *string {
						s := "0-1"
						return &s
					}(),
				},
				{
					GameID: 2,
					FEN: func() *string {
						s := "4Q2k/p5pp/3b3q/1p1p4/8/1B6/PPP3PP/5RK1 b - - 0 24"
						return &s
					}(),
					Result: func() *string {
						s := "1-0"
						return &s
					}(),
				},
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

			_ = <-round.LogChan

			client, err := round.CreateRoundClient()
			if err != nil {
				panic(err)
			}
			i := 0
			for {
				select {
				case update, ok := <-client.updates:
					if !ok {
						return
					}
					want := tt.wantMsgs[i]
					if !reflect.DeepEqual(update, tt.wantMsgs[i]) {
						deref := func(s *string) string {
							if s == nil {
								return "<nil>"
							}
							return *s
						}
						t.Errorf("Error in message number %d. Want: { %v %v %v %v }. Received: { %v %v %v %v }",
							i, want.GameID, deref(want.PGN), deref(want.FEN), deref(want.Result), update.GameID, deref(update.PGN), deref(update.FEN), deref(update.Result))
					}
					i += 1
				case <-timeout.C:
					t.Error("Timedout before receiving wanted messages")
					return
				}
			}
		})
	}
}
