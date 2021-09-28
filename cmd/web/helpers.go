package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/justinas/nosurf"
	"go.uber.org/zap"
	"watchess.org/watchess/pkg/models"
)

// Send 500 to user and log error
func (app *application) serverError(w http.ResponseWriter, err error) {
	app.logger.WithOptions(zap.AddCallerSkip(1)).Error("Server error", zap.Error(err), zap.ByteString("stack", debug.Stack()))
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Send client error to user with specified status
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// Wrapper around app.clientError to send a 404
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// Write rendered html to http.ResponseWriter
func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}

	// First write to buffer to make sure no error happens
	// If everything is okay, write buffer content to w
	buf := new(bytes.Buffer)
	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}

	buf.WriteTo(w)
}

// Return the User struct of the current logged in user or nil if no user
func (app *application) authenticatedUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(contextKeyUser).(*models.User)
	if !ok {
		return nil
	}
	return user
}

// Populate templateData with default values
func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	// Search form is present on every page
	td.CSRFToken = nosurf.Token(r)
	td.AuthenticatedUser = app.authenticatedUser(r)
	return td
}
