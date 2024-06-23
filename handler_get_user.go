package main

import (
	"database/sql"
	"net/http"
	"strings"

	"golang.org/x/net/context"
)

func (cfg *apiConfig) handler_get_user(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	reqApiKey := strings.TrimPrefix(authHeader, "ApiKey ")
	if reqApiKey == "" {
		respondWithError(w, http.StatusBadRequest, "no apiKey")
		return
	}

	ctx := context.Background()
	user, err := cfg.DB.GetUserByKey(ctx, reqApiKey)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "user not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, user)

}
