package main

import (
	"fmt"
	"net/http"

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
	app.render(w, r, "create.page.tmpl", &templateData{})
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

	id, err := app.tournaments.Insert(form.Get("name"), form.Get("short-description"), form.Get("long-description"), form.Has("standings"), *startDate, *endDate, false)

	if err != nil {
		app.serverError(w, err)
	}

	w.Write([]byte(fmt.Sprintf("Tournament created with id %v", id)))
}
