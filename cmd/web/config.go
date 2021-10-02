package main

import "watchess.org/watchess/pkg/models"

type config struct {
	tournament tournamentConfig
	user       userConfig
	round      roundConfig
}

type roundConfig struct {
	nameMax   int
	rewardMax int
	rewardMin int
}

type tournamentConfig struct {
	nameMax             int
	shortDescriptionMax int
	longDescriptionMax  int
}

type userConfig struct {
	nameMax     int
	emailMax    int
	passwordMin int
	validRoles  []string
}

func getConfig() config {
	return config{
		tournament: tournamentConfig{
			nameMax:             100,
			shortDescriptionMax: 400,
			longDescriptionMax:  20000,
		},
		user: userConfig{
			nameMax:     50,
			emailMax:    255,
			passwordMin: 8,
			validRoles:  models.GetUserRoleStrings(),
		},
		round: roundConfig{
			nameMax:   10,
			rewardMax: 100,
			rewardMin: 0,
		},
	}
}
