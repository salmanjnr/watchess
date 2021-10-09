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
	Win  float32 `json:"win"`
	Draw float32 `json:"draw"`
	Loss float32 `json:"loss"`
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
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	ShortDescription string    `json:"shortDescription"`
	LongDescription  string    `json:"longDescription"`
	HasStandings     bool      `json:"hasStandings"`
	StartDate        time.Time `json:"startDate"`
	EndDate          time.Time `json:"endDate"`
	IsLive           bool      `json:"isLive"`
	OwnerID          int       `json:"-"`
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
	ID           int        `json:"id"`
	Name         string     `json:"name"`
	PGNSource    string     `json:"-"`
	WhiteReward  GameReward `json:"whiteReward"`
	BlackReward  GameReward `json:"blackReward"`
	StartDate    time.Time  `json:"startDate"`
	TournamentID int        `json:"tournamentID"`
}

type Match struct {
	ID      int    `json:"id"`
	Side1   string `json:"side1"`
	Side2   string `json:"side2"`
	RoundID int    `json:"roundID"`
}

type Game struct {
	ID     int         `json:"id"`
	White  string      `json:"white"`
	Black  string      `json:"black"`
	Result *GameResult `json:"result"`
	// The side at which white player in Game model will be matched against in Match model
	// In case of a normal match this will just be player's name
	// In case of a team match this will be team name
	WhiteMatchSide string `json:"whiteMatchSide"`
	// The side at which black player in Game model will be matched against in Match model
	// In case of a normal match this will just be player's name
	// In case of a team match this will be team name
	BlackMatchSide string `json:"blackMatchSide"`
	PGN            string `json:"pgn"`
	MatchID        int    `json:"matchID"`
	RoundID        int    `json:"roundID"`
}
