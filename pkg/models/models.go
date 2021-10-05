package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
	ErrDuplicateUsername  = errors.New("models: duplicate username")
)

type UserRole uint8
type GameResult uint8

type GameReward struct {
	Win  float32
	Draw float32
	Loss float32
}

// Wrap roles in a function to prevent mutation
func GetUserRoleStrings() []string {
	return []string{"user", "admin"}
}

func GetGameResultStrings() []string {
	return []string{"1-0", "0.5-0.5", "0-1"}
}

const (
	NormalUser UserRole = iota
	AdminUser
)

const (
	WhiteWin GameResult = iota
	Draw
	BlackWin
)

func GetUserRole(roleStr string) (*UserRole, error) {
	userRoleStrings := GetUserRoleStrings()
	for i, v := range userRoleStrings {
		if v == roleStr {
			ui := UserRole(i)
			return &ui, nil
		}
	}
	return nil, errors.New("Role invalid")
}

func GetGameResult(resultStr string) (*GameResult, error) {
	gameResultStrings := GetGameResultStrings()
	for i, v := range gameResultStrings {
		if v == resultStr {
			ui := GameResult(i)
			return &ui, nil
		}
	}
	return nil, errors.New("Result invalid")
}

func (role UserRole) String() string {
	if role > 2 {
		return "Unknown"
	}
	return GetUserRoleStrings()[role]
}

func (res GameResult) String() string {
	if res > 2 {
		return "Unknown"
	}
	return GetGameResultStrings()[res]
}

type Tournament struct {
	ID               int
	Name             string
	ShortDescription string
	LongDescription  string
	HasStandings     bool
	StartDate        time.Time
	EndDate          time.Time
	IsLive           bool
	OwnerID          int
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Role           UserRole
}

type Round struct {
	ID           int
	Name         string
	PGNSource    string
	WhiteReward  GameReward
	BlackReward  GameReward
	StartDate    time.Time
	TournamentID int
}

type Match struct {
	ID      int
	Side1   string
	Side2   string
	RoundID int
}

type Game struct {
	ID             int
	White          string
	Black          string
	Result         *GameResult
	// The side at which white player in Game model will be matched against in Match model
	// In case of a normal match this will just be player's name 
	// In case of a team match this will be team name
	WhiteMatchSide string
	// The side at which black player in Game model will be matched against in Match model
	// In case of a normal match this will just be player's name 
	// In case of a team match this will be team name
	BlackMatchSide string
	MatchID        int
	RoundID        int
}
