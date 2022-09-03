package main

import (
	"Concurrency/section-6/final-project/constants"
	"net/http"
)

func (app *Config) SessionLoad(next http.Handler) http.Handler {
	return app.Session.LoadAndSave(next)
}

func (app *Config) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.Session.Exists(r.Context(), constants.UserIdTag) {
			app.Session.Put(r.Context(), constants.ErrorTag, "Log in first!")
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		next.ServeHTTP(w, r)
	})
}
