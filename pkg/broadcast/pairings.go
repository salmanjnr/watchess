// Provide order-agnostic match pairing
package broadcast

import (
	"errors"
	"sort"
	"sync"

	"github.com/mitchellh/hashstructure/v2"
	"github.com/salman69e27/chess"
)

type pairingHash uint64

type pairingValue []string

type pairing struct {
	Value pairingValue `hash:"set"`
}


func newPairing(side1, side2 string) pairing {
	return pairing {
		Value: pairingValue([]string{side1, side2}),
	}
}

func (p pairing) hash() pairingHash {
	// This should never fail
	hash, _ := hashstructure.Hash(p, hashstructure.FormatV2, nil)
	return pairingHash(hash)
}

func (p pairing) getSides() []string {
	s := []string(p.Value)
	sort.Strings(s)
	return s
}

// Get white and black players mapped to their corresponding match sides
func (p pairing) getPlayerSideMap() map[string]string {
	// Currently pairings are based on white and black tags, and are identical to game sides
	s := p.getSides()
	return map[string]string {
		s[0]: s[0],
		s[1]: s[1],
	}
}

// Map from pairing to id
type safePairingMap struct {
	v map[pairingHash]int
	mu sync.Mutex
}

func newPairingMap() safePairingMap {
	return safePairingMap{
		v: make(map[pairingHash]int),
	}
}

func (m *safePairingMap) value(p pairing) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	h := p.hash()
	if v, ok := m.v[h]; ok {
		return v, nil
	}
	return 0, errors.New("Pairing not present")
}

func (m *safePairingMap) create(p pairing, matchID int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	h := p.hash()
	if _, ok := m.v[h]; ok {
		return errors.New("Game already registered")
	}
	m.v[h] = matchID
	return nil
}

func (m *safePairingMap) delete(p pairing) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.v, p.hash())
}

// Helper function to get a match pairing from game tags. An error is returned if pairing can't be detected (this will only happen if White and Black tags are not present)
// Currently this always returns a pairing based on white and black tags, but will have methods for detecting team names in the future
func getPairingFromGame(gm *chess.Game) (pairing, error) {
	side1 := gm.GetTagPair("White")
	side2 := gm.GetTagPair("Black")
	if side1 == nil || side2 == nil {
		return pairing{}, errors.New("White and/or Black tags are missing")
	}
	return newPairing(side1.Value, side2.Value), nil
}
