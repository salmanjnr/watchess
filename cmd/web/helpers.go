package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
)

// Send 500 to user and log error
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
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
	err := ts.Execute(buf, td)
	if err != nil {
		app.serverError(w, err)
		return
	}

	buf.WriteTo(w)
}
