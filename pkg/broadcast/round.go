// Round struct is responsible for fetching PGN file from online source, detecting changes in the file, and sending changes to specific game centers.
package broadcast

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/notnil/chess"
)

type Fetcher struct {
	pgnSource url.URL
}

// Fetch from pgn source and return a reader. If any error happens, nil (along with the error) will be returned
func (f *Fetcher) Fetch() (io.Reader, error) {
	res, err := http.Get(f.pgnSource.String())
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

type Round struct {
	ID int
	// Maximum time without updates in pgn source before the round closes itself
	IdleTimeout time.Duration
	// The interval between each two successive updates
	UpdateInterval time.Duration
	// The fetcher that provides the pgn file
	pgnFetcher Fetcher
	// All registered game centers in the round. This is used for registering new clients instead of searching for the game id in Round.gms
	gBrokerMap safeGBrokerMap
	// Games ordered by their number in the pgn file
	gms         []*game
	roundBroker *roundBroker
}

func NewRound(roundID int, idleTimeout, updateInterval time.Duration, pgnSource url.URL) *Round {
	return &Round{
		ID:             roundID,
		IdleTimeout:    time.Duration(30) * time.Minute,
		UpdateInterval: time.Duration(500) * time.Millisecond,
		pgnFetcher:     Fetcher{pgnSource: pgnSource},
		gBrokerMap:     newGBrokerMap(),
		gms:            []*game{},
		roundBroker:    &roundBroker{newBroker()},
	}

}

// Creates game object and a record for the game center in Round.gameCenters
func (r *Round) createGame(gm *chess.Game, isFinished bool) error {
	id := rand.Int()

	// Only create a game center if the game is not already finished
	var gBroker *gameBroker
	if !isFinished {
		gBroker = &gameBroker{newBroker()}
		err := r.gBrokerMap.create(id, gBroker)
		if err != nil {
			return err
		}
		go gBroker.run()
	}

	r.gms = append(r.gms, &game{
		id:     id,
		chess:  gm,
		broker: gBroker,
	})
	return nil
}

// Creates a client for a game in the round. If game doesn't exist in game centers map (which either means the game is finished or never existed), an error is returned
func (r *Round) createGameClient(gameID int) (*GameClient, error) {
	gc, err := r.gBrokerMap.value(gameID)
	if err != nil {
		return nil, err
	}

	c := &GameClient{
		broker:  gc,
		updates: make(chan GameUpdate),
	}

	gc.registerChan <- c
	return c, nil
}

// Initialize round games. Creates a game center for all non-finished games and appends a game object in Rounds.gms
// An error is returned if one happened during parsing, including if a necessary game tag is not present
// Necessary tags: White, Black, Result
func (r *Round) Init() error {
	necessaryTags := []string{"White", "Black", "Result"}
	pgnReader, err := r.pgnFetcher.Fetch()
	if err != nil {
		return err
	}

	scanner := chess.NewScanner(pgnReader)

	for scanner.Scan() {
		game := scanner.Next()
		for _, tag := range necessaryTags {
			v := game.GetTagPair(tag)
			if v == nil {
				return fmt.Errorf("Game %v lacks a necessary tag: %v", len(r.gms)+1, tag)
			}
		}
		res := game.GetTagPair("Result")
		err := r.createGame(game, res.Value != "*")
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Round) handleUpdate(gameIndex int, currentGame *game, newGm *chess.Game) error {
	if currentGame.broker == nil {
		return fmt.Errorf("Cannot receive update for finished game number %v", gameIndex)
	}
	res := newGm.GetTagPair("Result")

	// Make sure result tag is valid
	if res.Value != "*" && res.Value != "1-0" && res.Value != "0-1" {
		return fmt.Errorf("Result value invalid for game %v", gameIndex)
	}

	pgnString := newGm.String()
	fen := newGm.FEN()
	if res.Value != "*" {
		// if game finished, send result and remove game center
		currentGame.broker.updates <- GameUpdate{GameID: currentGame.id, PGN: &pgnString, Result: &res.Value}
		r.roundBroker.updates <- GameUpdate{GameID: currentGame.id, FEN: &fen, Result: &res.Value}
		//r.gameCenters.delete(thisGame.center.ID)
		currentGame.broker.close()
		// We can't nullify the game center because it might still be running in another thread, so we just remove our reference to it and let the garbage collector do the heavy lifting
		*currentGame = game{}
	} else {
		currentGame.broker.updates <- GameUpdate{GameID: currentGame.id, PGN: &pgnString}
		r.roundBroker.updates <- GameUpdate{GameID: currentGame.id, FEN: &fen}
	}

	currentGame.chess = newGm
	return nil
}

func (r *Round) Close() {
	r.roundBroker.close()
}

func (r *Round) Run() error {
	defer r.Close()
	if r.roundBroker == nil {
		return fmt.Errorf("Round broker not initialized")
	}

	err := r.Init()
	if err != nil {
		return err
	}

	go r.roundBroker.run()

	timeoutTicker := time.NewTicker(r.IdleTimeout)
	defer timeoutTicker.Stop()

	updateTicker := time.NewTicker(r.UpdateInterval)
	defer updateTicker.Stop()

	for {
		select {
		case <-timeoutTicker.C:
			r.Close()
			break
		case <-updateTicker.C:
			pgnReader, err := r.pgnFetcher.Fetch()
			if err != nil {
				return err
			}
			scanner := chess.NewScanner(pgnReader)
			i := 0
			for scanner.Scan() {
				gm := scanner.Next()
				thisGame := r.gms[i]
				pgnString := gm.String()
				if thisGame.chess.String() != pgnString {
					err = r.handleUpdate(i, thisGame, gm)
					if err != nil {
						return err
					}
					// Reset ticker when a change is detected
					timeoutTicker.Reset(r.IdleTimeout)
				}
				i += 1
			}
		}
	}
}
