package broadcast

import (
	"bytes"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"watchess.org/watchess/pkg/models"
)

type fakeGameModel struct{}

func (m fakeGameModel) Insert(a, b string, c *models.GameResult, d, e, f string, g, h int) (int, error) {
	return rand.Int(), nil
}

type fakeMatchModel struct{}

func (m fakeMatchModel) Insert(a, b string, c int) (int, error) {
	return rand.Int(), nil
}

type testPGNServer struct {
	*httptest.Server
}

func (ts *testPGNServer) fetch() (io.Reader, error) {
	res, err := ts.Client().Get(ts.URL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}
	return bytes.NewReader(body), nil
}

func pgnHandlerFactory(pgns []string) func(w http.ResponseWriter, r *http.Request) {
	timeCreated := time.Now()
	sort.Strings(pgns)
	return func(w http.ResponseWriter, r *http.Request) {
		diff := int(time.Now().Sub(timeCreated).Milliseconds() / 50)
		if diff >= len(pgns) {
			diff = len(pgns) - 1
		}

		content, err := ioutil.ReadFile(pgns[diff])
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(content))
	}
}

func newPGNServer(t *testing.T, dir string) (*testPGNServer, error) {
	pgns, err := filepath.Glob(filepath.Join("testdata", dir, "*.pgn"))
	if err != nil {
		return nil, err
	}
	h := pgnHandlerFactory(pgns)
	ts := httptest.NewServer(http.HandlerFunc(h))
	return &testPGNServer{ts}, nil
}

func newTestRound(ts *testPGNServer) *Round {
	return &Round{
		roundID:        0,
		idleTimeout:    time.Duration(500) * time.Millisecond,
		updateInterval: time.Duration(10) * time.Millisecond,
		pgnFetcher:     ts,
		gBrokerMap:     newGBrokerMap(),
		gms:            []*game{},
		matchModel:     fakeMatchModel{},
		gameModel:      fakeGameModel{},
		roundBroker:    &roundBroker{newBroker()},
		pairingMap:     newPairingMap(),
		errorChan:      make(chan error),
		LogChan:        make(chan interface{}),
		Done:           make(chan struct{}),
	}
}
