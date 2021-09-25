package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Tournament struct {
	ID               int
	Name             string
	ShortDescription string
	LongDescription  string
	HasStandings     bool
	StartDate        time.Time
	EndDate          time.Time
	IsLive           bool
}
