package main

import (
	"html/template"
	"path/filepath"

	"watchess.org/watchess/pkg/forms"
	"watchess.org/watchess/pkg/models"
)

type templateData struct {
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
