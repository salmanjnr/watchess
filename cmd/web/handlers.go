package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"watchess.org/watchess/pkg/forms"
	"watchess.org/watchess/pkg/models"
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

	if len(active)+len(finished)+len(upcoming) != 0 {
		td.Tournaments = &tournaments{
			Active:   active,
			Finished: finished,
			Upcoming: upcoming,
		}
	}

	app.render(w, r, "home.page.tmpl", td)
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	if app.authenticatedUser(r) != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	// If user already authenticated, direct them to homepage

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	form := forms.New(r.PostForm)
	form.Set("role", "admin") // currently, there is no normal users

	form.Required("name", "email", "password", "confirm-password")
	form.ValidEmail("email")
	form.Matching("password", "confirm-password")
	form.MinLength("password", app.config.user.passwordMin)
	form.MaxLength("name", app.config.user.nameMax)
	form.MaxLength("email", app.config.user.emailMax)
	// This will never generate an error because we just set the role value, but we keep it for the future
	form.PermittedValues("role", app.config.user.validRoles...)

	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{
			Form: form,
		})
		return
	}

	userRole, err := models.GetUserRole(form.Get("role"))
	// At this point, no error should happen because we just validated that the role field has a valid value
	if err != nil {
		app.serverError(w, err)
	}

	id, err := app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"), *userRole)

	if err == models.ErrDuplicateEmail {
		form.Errors.Add("email", "Email already registered")
		app.render(w, r, "signup.page.tmpl", &templateData{
			Form: form,
		})
		return
	} else if err == models.ErrDuplicateUsername {
		form.Errors.Add("name", "Username taken")
		app.render(w, r, "signup.page.tmpl", &templateData{
			Form: form,
		})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "userID", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err == models.ErrInvalidCredentials {
		form.Errors.Add("generic", "Email or Password is incorrect")
		app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "userID", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "userID")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) createTournamentForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create-tournament.page.tmpl", &templateData{})
}

func (app *application) createTournament(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
	}

	form := forms.New(r.PostForm)
	form.Required("name", "short-description", "start-date", "end-date")
	form.MaxLength("name", app.config.tournament.nameMax)
	form.MaxLength("short-description", app.config.tournament.shortDescriptionMax)
	form.MaxLength("long-description", app.config.tournament.longDescriptionMax)
	startDate, endDate := form.DatePair("start-date", "end-date")

	if !form.Valid() {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	user := app.authenticatedUser(r)
	if user == nil {
		app.serverError(w, errors.New("User is absent from context in a handler that require auth"))
	}

	id, err := app.tournaments.Insert(form.Get("name"), form.Get("short-description"), form.Get("long-description"), form.Has("standings"), *startDate, *endDate, false, user.ID)

	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Write([]byte(fmt.Sprintf("Tournament created with id %v", id)))
}

func (app *application) createRoundForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	tournament, err := app.tournaments.Get(id)

	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "create-round.page.tmpl", &templateData{Tournament: tournament})
}

func (app *application) createRound(w http.ResponseWriter, r *http.Request) {
	tournamentId, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || tournamentId < 1 {
		app.notFound(w)
		return
	}

	tournament, err := app.tournaments.Get(tournamentId)

	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		app.serverError(w, err)
	}

	form := forms.New(r.PostForm)
	form.Required("name", "pgn-source", "start-date")
	form.MaxLength("name", app.config.round.nameMax)
	form.ValidURL("pgn-source")
	rewardFields := []string{"white-win", "white-draw", "white-loss", "black-win", "black-draw", "black-loss"}
	form.Numerical(rewardFields...)
	for _, field := range rewardFields {
		form.MaxValue(field, app.config.round.rewardMax)
		form.MinValue(field, app.config.round.rewardMin)
	}
	date := form.Date("start-date")

	app.logger.Debug("Round Creation", zap.String("form", fmt.Sprintf("%v", form)), zap.String("date", date.String()))
	if !form.Valid() {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	atoiWithDefault := func (str string, defaultValue float32) float32 {
		str = strings.TrimSpace(str)
		if str == "" {
			return defaultValue
		}
		nbr, err := strconv.ParseFloat(str, 32)
		if err != nil {
			return defaultValue
		}
		return float32(nbr)
	}

	whiteReward := models.GameReward{
		Win: atoiWithDefault(form.Get("white-win"), 1),
		Draw: atoiWithDefault(form.Get("white-draw"), 0.5),
		Loss: atoiWithDefault(form.Get("white-loss"), 0),
	}

	blackReward := models.GameReward{
		Win: atoiWithDefault(form.Get("black-win"), 1),
		Draw: atoiWithDefault(form.Get("black-draw"), 0.5),
		Loss: atoiWithDefault(form.Get("black-loss"), 0),
	}

	id, err := app.rounds.Insert(form.Get("name"), form.Get("pgn-source"), whiteReward, blackReward, *date, tournament.ID)

	if err != nil {
		app.serverError(w, err)
		return
	}
	round, err := app.rounds.Get(id)
	app.logger.Debug("Round Inserted", zap.Error(err))
	w.Write([]byte(fmt.Sprintf("%v", round)))
}
