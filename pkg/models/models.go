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
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Role           UserRole
}
