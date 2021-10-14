package broadcast

import (
	"errors"
	"fmt"
	"sync"

	"github.com/notnil/chess"
)

// All communication is done in terms of GameUpdate
type GameUpdate struct {
	GameID int
	// PGN is used only when sending an update via a game broker
	PGN *string
	// FEN is used only when sending an update via a round broker
	FEN *string
	// Result can be used with a game broker or a round broker
	Result *string
}

type safeGBrokerMap struct {
	v  map[int]*gameBroker
	mu sync.Mutex
}

type game struct {
	id     int
	chess  *chess.Game
	broker *gameBroker
}

func newGBrokerMap() safeGBrokerMap {
	return safeGBrokerMap{
		v: make(map[int]*gameBroker),
	}
}

func (g *safeGBrokerMap) value(gameID int) (*gameBroker, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if v, ok := g.v[gameID]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("Game center with ID %v not present", gameID)
}

func (g *safeGBrokerMap) create(gameID int, gc *gameBroker) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.v[gameID]; ok {
		return errors.New("Game already registered")
	}
	g.v[gameID] = gc
	return nil
}

func (g *safeGBrokerMap) delete(gameID int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.v, gameID)
}
