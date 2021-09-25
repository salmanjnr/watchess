package main

type config struct {
	tournament tournamentConfig
}

type tournamentConfig struct {
	nameMax             int
	shortDescriptionMax int
	longDescriptionMax  int
}

func getConfig() config {
	return config{
		tournament: tournamentConfig{
			nameMax:             100,
			shortDescriptionMax: 400,
			longDescriptionMax:  20000,
		},
	}
}
