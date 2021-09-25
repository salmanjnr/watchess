package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New()
	dynamicMiddleware := alice.New(app.session.Enable)
	formMiddleware := dynamicMiddleware.Append(noSurf)

	mux := pat.New()

	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/tournaments/create", formMiddleware.ThenFunc(app.createTournamentForm))
	mux.Post("/tournaments/create", formMiddleware.ThenFunc(app.createTournament))
	mux.Get("/ping", http.HandlerFunc(ping))

	fileServer := http.FileServer(http.Dir("./ui/"))
	mux.Get("/ui/", http.StripPrefix("/ui", fileServer))

	return standardMiddleware.Then(mux)
}
