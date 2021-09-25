package main

import (
	"net/http"

	"github.com/bmizerany/pat"
)

func (app *application) routes() http.Handler {
	mux := pat.New()

	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/ping", http.HandlerFunc(ping))
	mux.Get("/tournaments/create", http.HandlerFunc(app.createTournamentForm))
	mux.Post("/tournaments/create", http.HandlerFunc(app.createTournament))

	fileServer := http.FileServer(http.Dir("./ui/"))
	mux.Get("/ui/", http.StripPrefix("/ui", fileServer))

	return mux
}
