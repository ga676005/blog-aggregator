package main

import "net/http"

type readiness struct {
	Status string `json:"status"`
}

func handler_readiness(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, readiness{"ok"})
}

func handler_error(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}
