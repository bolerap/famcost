package main

import "net/http"

func checkInternalServerError(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func isAuthenticated(w http.ResponseWriter, r *http.Request) {
	if !authenticated {
		http.Redirect(w, r, "/login", 301)
	}
}
