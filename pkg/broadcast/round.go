// Round struct is responsible for fetching PGN file from online source, detecting changes in the file, and sending changes to specific brokers.
package broadcast

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/salman69e27/chess"
	"watchess.org/watchess/pkg/models"
	"watchess.org/watchess/pkg/models/mysql"
)

var requiredTags = []string{"White", "Black", "Result"}

type fetcher interface {
	fetch() (io.Reader, error)
}

type pgnFetcher struct {
	pgnSource url.URL
}

// fetch from pgn source and return a reader. If any error happens, nil (along with the error) will be returned
func (f pgnFetcher) fetch() (io.Reader, error) {
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
	roundID int
	// Maximum time without updates in pgn source before the round closes itself
	idleTimeout time.Duration
	// The interval between each two successive updates
	updateInterval time.Duration
	// The fetcher that provides the pgn file
	pgnFetcher fetcher
	// All registered game brokers in the round. This is used for registering new clients instead of searching for the game id in Round.gms
	gBrokerMap safeGBrokerMap
	// Games ordered by their number in the pgn file
	gms        []*game
	matchModel interface {
		Insert(string, string, int) (int, error)
	}
	gameModel interface {
		Insert(string, string, *models.GameResult, string, string, string, int, int) (int, error)
	}
	roundBroker *roundBroker
	pairingMap  safePairingMap
	// Channel used by brokers for reporting errors to round
	errorChan chan error
	// Channel used by round to report updates summary
	LogChan chan interface{}
	// Channel to close the round
	Done chan struct{}
}

func NewRound(roundID int, idleTimeout, updateInterval time.Duration, pgnSource url.URL, matchModel *mysql.MatchModel, gameModel *mysql.GameModel) *Round {
	return &Round{
		roundID:        roundID,
		idleTimeout:    idleTimeout,
		updateInterval: updateInterval,
		pgnFetcher:     pgnFetcher{pgnSource: pgnSource},
		gBrokerMap:     newGBrokerMap(),
		gms:            []*game{},
		matchModel:     matchModel,
		gameModel:      gameModel,
		roundBroker:    &roundBroker{newBroker()},
		pairingMap:     newPairingMap(),
		errorChan:      make(chan error),
		LogChan:        make(chan interface{}),
		Done:           make(chan struct{}),
	}

}

// Creates game (and match if necessary) record in database, game object and broker. Returns (isMatchCreated, isGameCreated, error)
func (r *Round) createGame(gm *chess.Game) (matchCreated bool, gameCreated bool, err error) {
	for _, tag := range requiredTags {
		v := gm.GetTagPair(tag)
		if v == nil {
			err = fmt.Errorf("Game %d missing required tag: %s", len(r.gms), tag)
			return
		}
	}
	res := gm.GetTagPair("Result").Value
	white := gm.GetTagPair("White").Value
	black := gm.GetTagPair("Black").Value

	p, err := getPairingFromGame(gm)
	if err != nil {
		return
	}

	matchID, err := r.pairingMap.value(p)
	if err != nil {
		sides := p.getSides()
		matchID, err = r.matchModel.Insert(sides[0], sides[1], r.roundID)
		if err != nil {
			return
		}
		matchCreated = true
		err = r.pairingMap.create(p, matchID)
		if err != nil {
			return
		}
	}

	var gameRes *models.GameResult = nil
	if res != "*" {
		gameRes, err = models.GetGameResult(res)
		if err != nil {
			return
		}
	}

	playerSides := p.getPlayerSideMap()

	gameID, err := r.gameModel.Insert(white, black, gameRes, playerSides[white], playerSides[black], gm.String(), matchID, r.roundID)

	if err != nil {
		return
	}
	gameCreated = true
	// Only create a game broker if the game is not already finished
	var gBroker *gameBroker
	if gameRes == nil {
		gBroker = &gameBroker{newBroker()}
		err = r.gBrokerMap.create(gameID, gBroker)
		if err != nil {
			return
		}
		go gBroker.run(r.errorChan)
	}

	r.gms = append(r.gms, &game{
		id:     gameID,
		chess:  gm,
		broker: gBroker,
	})
	return
}

// Creates a client for a game in the round. If game doesn't exist in game brokers map (which either means the game is finished or never existed), an error is returned
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

// Initialize round games. Creates a game broker for all non-finished games and appends a game object in Rounds.gms
// An error is returned if one happened during parsing, including if a necessary game tag is not present
// Necessary tags: White, Black, Result
func (r *Round) Init() error {
	go r.roundBroker.run(r.errorChan)
	pgnReader, err := r.pgnFetcher.fetch()
	if err != nil {
		return err
	}

	scanner := chess.NewScanner(pgnReader)
	numMatches := 0
	numGames := 0
	for i := 0; scanner.Scan(); i += 1 {
		gm := scanner.Next()
		matchCreated, gameCreated, err := r.createGame(gm)
		if err != nil {
			return err
		}
		if matchCreated {
			numMatches += 1
		}
		if gameCreated {
			numGames += 1
		}
	}
	if scanner.Err() != io.EOF {
		return scanner.Err()
	}

	r.LogChan <- fmt.Sprintf("%d matches and %d games detected", numMatches, numGames)
	return nil
}

func (r *Round) handleUpdate(gameIndex int, currentGame *game, newGm *chess.Game) error {
	if currentGame.broker == nil {
		return fmt.Errorf("Cannot receive update for finished game number %v", gameIndex)
	}
	res := newGm.GetTagPair("Result")

	// Make sure result tag is valid
	if res.Value != "*" && res.Value != "1-0" && res.Value != "0-1" && res.Value != "1/2-1/2" {
		return fmt.Errorf("Result value %s invalid for game %v", res.Value, gameIndex)
	}

	pgnString := newGm.String()
	fen := newGm.FEN()
	if res.Value != "*" {
		// if game finished, send result and remove game broker
		currentGame.broker.updates <- GameUpdate{GameID: currentGame.id, PGN: &pgnString, Result: &res.Value}
		r.roundBroker.updates <- GameUpdate{GameID: currentGame.id, FEN: &fen, Result: &res.Value}
		r.gBrokerMap.delete(currentGame.id)
		currentGame.broker.close()
		// We can't nullify the game broker because it might still be running in another thread, so we just remove our reference to it and let the garbage collector do the heavy lifting
		*currentGame = game{}
	} else {
		currentGame.broker.updates <- GameUpdate{GameID: currentGame.id, PGN: &pgnString}
		r.roundBroker.updates <- GameUpdate{GameID: currentGame.id, FEN: &fen}
	}

	currentGame.chess = newGm
	return nil
}

func (r *Round) Close() {
	close(r.LogChan)
	for _, gm := range r.gms {
		if gm.broker != nil {
			gm.broker.close()
		}
	}
	r.roundBroker.close()
	// We can't close the error channel because brokers might be using it, so we just block them
	// and they will close with the termination signal
	r.errorChan = nil
}

func (r *Round) Run() {
	defer func() {
		r.LogChan <- "Closing"
		r.Close()
	}()
	if r.roundBroker == nil {
		r.LogChan <- fmt.Errorf("Round broker not initialized")
		return
	}

	err := r.Init()
	if err != nil {
		r.LogChan <- err
		return
	}

	timeoutTicker := time.NewTicker(r.idleTimeout)
	defer timeoutTicker.Stop()

	updateTicker := time.NewTicker(r.updateInterval)
	defer updateTicker.Stop()

	for {
		select {
		case <-r.Done:
			return
		case err = <-r.errorChan:
			r.LogChan <- err
			return
		case <-timeoutTicker.C:
			return
		case <-updateTicker.C:
			pgnReader, err := r.pgnFetcher.fetch()
			if err != nil {
				r.LogChan <- err
				return
			}
			scanner := chess.NewScanner(pgnReader)
			updatedGames := 0
			allGames := ""
			for i := 0; scanner.Scan(); i += 1 {
				gm := scanner.Next()
				allGames += gm.String() + "\n"
				thisGame := r.gms[i]
				pgnString := gm.String()
				if thisGame.chess.String() != pgnString {
					err = r.handleUpdate(i, thisGame, gm)
					if err != nil {
						r.LogChan <- err
						return
					}
					// Reset ticker when a change is detected
					timeoutTicker.Reset(r.idleTimeout)
					updatedGames += 1
				}
			}
			if updatedGames != 0 {
				if updatedGames == 1 {
					r.LogChan <- "1 game updated"
				} else {
					r.LogChan <- fmt.Sprintf("%d games updated", updatedGames)
				}
			}
		}
	}
}
