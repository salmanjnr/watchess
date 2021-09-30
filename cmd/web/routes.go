package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New()
	dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)
	authDynamicMiddleware := dynamicMiddleware.Append(app.requireAdminUser)

	mux := pat.New()

	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))

	mux.Get("/tournaments/create", authDynamicMiddleware.ThenFunc(app.createTournamentForm))
	mux.Post("/tournaments/create", authDynamicMiddleware.ThenFunc(app.createTournament))

	mux.Get("/tournaments/:id/rounds/create", authDynamicMiddleware.ThenFunc(app.createRoundForm))
	mux.Post("/tournaments/:id/rounds/create", authDynamicMiddleware.ThenFunc(app.createRound))

	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", authDynamicMiddleware.ThenFunc(app.logoutUser))

	mux.Get("/ping", http.HandlerFunc(ping))

	fileServer := http.FileServer(http.Dir("./ui/"))
	mux.Get("/ui/", http.StripPrefix("/ui", fileServer))

	return standardMiddleware.Then(mux)
}
