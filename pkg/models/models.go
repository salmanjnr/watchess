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

type GameReward struct {
	Win  float32
	Draw float32
	Loss float32
}

// Wrap roles in a function to prevent mutation
func GetUserRoleStrings() []string {
	return []string{"user", "admin"}
}

const (
	NormalUser UserRole = iota
	AdminUser
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

func (role UserRole) String() string {
	if role > 2 {
		return "Unknown"
	}
	return GetUserRoleStrings()[role]
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
