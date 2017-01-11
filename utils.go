package main

import "net/http"

func checkInternalServerError(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func checkConflict(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, "The user was existed", http.StatusConflict)
	}
}
