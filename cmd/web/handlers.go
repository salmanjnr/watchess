package main

import (
	"fmt"
	"net/http"

	"watchess.org/watchess/pkg/forms"
)

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	td := &templateData{}

	active, err := app.tournaments.LatestActive(5)
	if err != nil {
		app.serverError(w, err)
	}

	finished, err := app.tournaments.LatestFinished(3)
	if err != nil {
		app.serverError(w, err)
	}

	upcoming, err := app.tournaments.Upcoming(3)
	if err != nil {
		app.serverError(w, err)
	}

	if len(active) + len(finished) + len(upcoming) != 0 {
		td.Tournaments = &tournaments{
			Active: active,
			Finished: finished,
			Upcoming: upcoming,
		}
	}

	app.render(w, r, "home.page.tmpl", td)
}

func (app *application) createTournamentForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", nil)
}

func (app *application) createTournament(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
	}

	form := forms.New(r.PostForm)
	form.Required("name", "short-description", "long-description", "start-date", "end-date")
	form.MaxLength("short-description", 100)
	startDate, endDate := form.DatePair("start-date", "end-date")

	if !form.Valid() {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	id, err := app.tournaments.Insert(form.Get("name"), form.Get("short-description"), form.Get("long-description"), form.Has("standings"), *startDate, *endDate, false)

	if err != nil {
		app.serverError(w, err)
	}

	w.Write([]byte(fmt.Sprintf("Tournament created with id %v", id)))
}
