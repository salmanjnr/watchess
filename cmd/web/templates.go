package main

import (
	"errors"
	"html/template"
	"path/filepath"

	"watchess.org/watchess/pkg/forms"
	"watchess.org/watchess/pkg/models"
)

type templateData struct {
	Rounds            []*models.Round
	Tournament        *models.Tournament
	Tournaments       *tournaments
	CSRFToken         string
	Form              *forms.Form
	AuthenticatedUser *models.User
}

type tournaments struct {
	Active   []*models.Tournament
	Upcoming []*models.Tournament
	Finished []*models.Tournament
}

type roundGames struct {
	models.Round
	Matches      []matchGames      `json:"matches"`
}

type matchGames struct {
	models.Match
	Games   []models.Game `json:"games"`
}

func (app *application) newRoundGames(round *models.Round) (*roundGames, error) {
	if round == nil {
		return nil, errors.New("Round can't be nil")
	}
	matches := []matchGames{}

	ms, err := app.matches.GetByRound(round.ID)
	if err != nil {
		return nil, err
	}

	// Instead of iterating over matches and getting games by match, we 
	// get all round games and store them in a map with match id as key
	gs, err := app.games.GetByRound(round.ID)
	if err != nil {
		return nil, err
	}

	matchesMap := map[int]*[]models.Game{}

	for _, game := range gs {
		if p, ok := matchesMap[game.MatchID]; ok {
			*p = append(*p, *game)
		} else {
			matchesMap[game.MatchID] = &[]models.Game{*game}
		}
	}

	for _, match := range ms {
		var curGames []models.Game
		if games, ok := matchesMap[match.ID]; ok{
			curGames = *games
		}

		matches = append(matches, matchGames{
			Match: *match,
			Games: curGames,
		})
	}

	rgames := roundGames{
		Round: *round,
		Matches: matches,
	}
	return &rgames, nil
}

// Initialize cache for all template files in dir
func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// Load all page templates
	pages, err := filepath.Glob(filepath.Join(dir, "*.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// Template functions need to be registered before parsing template files
		ts, err := template.New(name).ParseFiles(page)

		if err != nil {
			return nil, err
		}

		// Parse all layouts and partials
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}

	return cache, nil
}
