package main

import (
	"net/http"

	"github.com/bmizerany/pat"
)

func (app *application) routes() http.Handler {
	mux := pat.New()
	// User routes
	mux.Post("/api/user", http.HandlerFunc(app.createUser))
	mux.Get("/api/user/:id", http.HandlerFunc(app.getUser))
	mux.Put("/api/user/:id", http.HandlerFunc(app.updateUser))
	mux.Get("/api/user", http.HandlerFunc(app.listUser))
	mux.Del("/api/user/:id", http.HandlerFunc(app.deleteUser))

	// Store routes
	// list stores
	// get store by ID
	mux.Post("/api/store", http.HandlerFunc(app.createStore))
	mux.Put("/api/store/:id", http.HandlerFunc(app.updateStore))
	// delete store

	return app.logRequest(mux)
}
